package managerscheduler

import (
	"context"
	"errors"
	"time"

	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	problemsrepo "github.com/keepcalmist/chat-service/internal/repositories/problems"
	managerpool "github.com/keepcalmist/chat-service/internal/services/manager-pool"
	"github.com/keepcalmist/chat-service/internal/types"
	"go.uber.org/zap"
)

const serviceName = "manager-scheduler"

type transactor interface {
	RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}

type managerPool interface {
	Get(ctx context.Context) (types.UserID, error)
	Put(ctx context.Context, managerID types.UserID) error
}

type problemsRepo interface {
	GetUnassignedProblems(ctx context.Context) ([]*problemsrepo.Problem, error)
	GetRequestID(ctx context.Context, problemID types.ProblemID) (types.RequestID, error)
	SetManagerForProblem(ctx context.Context, problemID types.ProblemID, managerID types.UserID) error
}

type outboxSvc interface {
	Put(ctx context.Context, name string, payload string, availableAt time.Time) (types.JobID, error)
}

type msgRepo interface {
	CreateServiceMessage(
		ctx context.Context,
		requestID types.RequestID,
		chatID types.ChatID,
		problemID types.ProblemID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

//go:generate options-gen -out-filename=service_options.gen.go -from-struct=Options
type Options struct {
	period time.Duration `option:"mandatory" validate:"min=100ms,max=1m"`

	managerPool  managerPool  `option:"mandatory" validate:"required"`
	msgRepo      msgRepo      `option:"mandatory" validate:"required"`
	problemsRepo problemsRepo `option:"mandatory" validate:"required"`
	outboxSvc    outboxSvc    `option:"mandatory" validate:"required"`
	transactor   transactor   `option:"mandatory" validate:"required"`
}

type Service struct {
	Options

	logger *zap.Logger
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &Service{
		Options: opts,
		logger:  zap.L().Named(serviceName),
	}, nil
}

func (s *Service) Run(ctx context.Context) error {
mainLoop:
	for {
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(s.period):
		}

		unassignedProblems, err := s.problemsRepo.GetUnassignedProblems(ctx)
		if err != nil {
			s.logger.Error("failed to get unassigned problems", zap.Error(err))
			continue
		}

		for _, problem := range unassignedProblems {
			managerID, err := s.managerPool.Get(ctx)
			if err != nil {
				if errors.Is(err, managerpool.ErrNoAvailableManagers) {
					continue mainLoop
				}

				s.logger.Error("failed to get manager from pool", zap.Error(err))
				continue
			}

			_, err = s.msgRepo.CreateServiceMessage(
				ctx,
				problem.RequestID,
				problem.ChatID,
				problem.ID,
				"Manager "+managerID.String()+" will answer you",
			)
			if err != nil {
				s.logger.Error("failed to create service message", zap.Error(err))
				continue
			}
		}
	}
	// FIXME: Каждый s.period времени необходимо:
	// FIXME: - достать все проблемы, ожидающие менеджера;
	// FIXME: - назначить на них менеджеров, ожидающих в пуле.

	// FIXME: Самой "старой" проблеме должен достаться дольше всех ожидающий менеджер.
	// FIXME: Иными словами – не должно быть ситуации, когда кто-то обратился в
	// FIXME: поддержку банка после тебя, но при этом раньше получил ответ.
	// FIXME: Также не должно быть ситуации, когда ты как менеджер заступил на смену
	// FIXME: и сразу же получил клиента, хотя твой коллега заступил на смену раньше
	// FIXME: и сидит без работы.

	// FIXME: Проблема ожидает менеджера, если:
	// FIXME: 1) нет назначенного на неё менеджера;
	// FIXME: 2) в чате есть сообщения, видимые менеджеру.

	// FIXME: Обратите внимание, что в общем случае количество
	// FIXME: проблем сильно больше ожидающих менеджеров!

	// FIXME: Если менеджер успешно назначен на проблему, то клиенту должно улететь служебное сообщение
	// FIXME: с текстом "Manager <id> will answer you".

	// FIXME: Служебное сообщение – сообщение, у которого флаг is_service == true и author_id == types.UserIDNil.

	return nil
}

func (s *Service) assignManagerToProblem(ctx context.Context, p *problemsrepo.Problem) error {
	managerID, err := s.managerPool.Get(ctx)
	if err != nil {
		if errors.Is(err, managerpool.ErrNoAvailableManagers) {
			return nil
		}
	}

	err = s.transactor.RunInTx(ctx, func(ctx context.Context) error {
		err := s.problemsRepo.SetManagerForProblem(ctx, p.ID, managerID)
		if err != nil {
			return err
		}

		reqID, err := s.problemsRepo.GetRequestID(ctx, p.ID)

		_, err = s.outboxSvc.Put(ctx, "manager-scheduler", problemID.String(), time.Now())
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		putErr := s.managerPool.Put(ctx, managerID)
		if putErr != nil {
			s.logger.Error("failed to put manager to pool", zap.Error(putErr))
		}

		return err
	}

	return nil
}
