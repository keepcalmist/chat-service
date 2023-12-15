package outbox

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.uber.org/zap"

	jobsrepo "github.com/keepcalmist/chat-service/internal/repositories/jobs"
	"github.com/keepcalmist/chat-service/internal/types"
)

const serviceName = "outbox"

var ErrJobAlreadyRegistered = errors.New("job already registered")

type jobsRepository interface {
	CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
	FindAndReserveJob(ctx context.Context, until time.Time) (jobsrepo.Job, error)
	CreateFailedJob(ctx context.Context, name, payload, reason string) error
	DeleteJob(ctx context.Context, jobID types.JobID) error
}

type transactor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	workers    int            `option:"mandatory" validate:"min=1,max=32"`
	idleTime   time.Duration  `option:"mandatory" validate:"min=100ms,max=10s"`
	reserveFor time.Duration  `option:"mandatory" validate:"min=1s,max=10m"`
	r          jobsRepository `option:"mandatory"`
	t          transactor     `option:"mandatory"`
	logger     *zap.Logger
}

type Service struct {
	Options
	jobs map[string]Job
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if opts.logger == nil {
		opts.logger = zap.L().Named(serviceName)
	}

	return &Service{
		Options: opts,
		jobs:    make(map[string]Job),
	}, nil
}

func (s *Service) RegisterJob(job Job) error {
	if _, ok := s.jobs[job.Name()]; ok {
		return fmt.Errorf("%w: %s", ErrJobAlreadyRegistered, job.Name())
	}
	s.jobs[job.Name()] = job
	return nil
}

func (s *Service) MustRegisterJob(job Job) {
	if err := s.RegisterJob(job); err != nil {
		panic(err)
	}
}

func (s *Service) Run(ctx context.Context) error {
	wg := new(sync.WaitGroup)
	for i := 0; i < s.workers; i++ {
		wg.Add(1)

		s.logger.Info("starting worker", zap.Int("worker_id", i))
		go func() {
			s.startWorker(ctx)
			defer wg.Done()
		}()
	}
	wg.Wait()

	return nil
}

func (s *Service) startWorker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			s.logger.Debug("context done")
			return
		default:
			reservedJob, err := s.r.FindAndReserveJob(ctx, time.Now().Add(s.reserveFor))
			if err != nil {
				if errors.Is(err, jobsrepo.ErrNoJobs) {
					s.logger.Debug("sleeping", zap.Duration("idle_time", s.idleTime))
					select {
					case <-ctx.Done():
						s.logger.Error("context done", zap.Error(ctx.Err()))
						return
					case <-time.After(s.idleTime):
					}
					continue
				}
				s.logger.Error("failed to find and reserve job", zap.Error(err))
				continue
			}

			s.logger.Info("job found, start processing...", zap.String("job_id", reservedJob.ID.String()))

			j, ok := s.jobs[reservedJob.Name]
			if !ok {
				err = s.CreateFailedAndDeleteMainJob(ctx, reservedJob)
				if err != nil {
					s.logger.Error("failed to create failed job", zap.Error(err),
						zap.String("job_id", reservedJob.ID.String()))
					continue
				}

				continue
			}

			if reservedJob.Attempts > j.MaxAttempts() {
				err = s.CreateFailedAndDeleteMainJob(ctx, reservedJob)
				if err != nil {
					s.logger.Error("failed to create failed job", zap.Error(err),
						zap.String("job_id", reservedJob.ID.String()))
					continue
				}
				continue
			}

			err = s.handleJob(ctx, reservedJob, j)
			if err != nil {
				s.logger.Error("failed to handle job", zap.Error(err), zap.String("job_id", reservedJob.ID.String()))
				if j.MaxAttempts() < reservedJob.Attempts {
					err = s.CreateFailedAndDeleteMainJob(ctx, reservedJob)
					if err != nil {
						s.logger.Error("failed to create failed job", zap.Error(err),
							zap.String("job_id", reservedJob.ID.String()))
						continue
					}

					continue
				}
			}
		}
	}
}

func (s *Service) handleJob(ctx context.Context, reservedJob jobsrepo.Job, j Job) error {
	ctxWithCancel, cancel := context.WithTimeout(ctx, j.ExecutionTimeout())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		errChan <- j.Handle(ctxWithCancel, reservedJob.Payload)
	}()

	select {
	case <-ctxWithCancel.Done():
		s.logger.Error("context done", zap.Error(ctxWithCancel.Err()))
		return ctxWithCancel.Err()
	case err := <-errChan:
		if err != nil {
			return err
		}
	}

	if err := s.r.DeleteJob(ctx, reservedJob.ID); err != nil {
		s.logger.With().Error("failed to delete job", zap.Error(err), zap.String("job_id", reservedJob.ID.String()))
		return err
	}

	s.logger.Info("job successfully handled", zap.String("job_id", reservedJob.ID.String()))

	return nil
}

func (s *Service) CreateFailedAndDeleteMainJob(ctx context.Context, job jobsrepo.Job) error {
	err := s.t.RunInTx(ctx, func(ctx context.Context) error {
		if err := s.r.CreateFailedJob(ctx, job.Name, job.Payload, "max attempts exceeded"); err != nil {
			return fmt.Errorf("failed to create failed job: %w", err)
		}

		if err := s.r.DeleteJob(ctx, job.ID); err != nil {
			return fmt.Errorf("failed to delete job: %w", err)
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to run transaction: %w", err)
	}

	s.logger.Info("job failed and deleted", zap.String("job_id", job.ID.String()))

	return nil
}
