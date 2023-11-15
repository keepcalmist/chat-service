package jobsrepo

import (
	"context"
	"errors"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/keepcalmist/chat-service/internal/store"
	"github.com/keepcalmist/chat-service/internal/store/job"
	"github.com/keepcalmist/chat-service/internal/types"
)

var ErrNoJobs = errors.New("no jobs found")

type Job struct {
	ID       types.JobID
	Name     string
	Payload  string
	Attempts int
}

func (r *Repo) FindAndReserveJob(ctx context.Context, until time.Time) (Job, error) {
	retJob := Job{}
	err := r.db.RunInTx(ctx, func(ctx context.Context) error {
		j, err := r.db.Job(ctx).Query().
			Unique(false).
			Where(
				job.And(
					job.AvailableAtLTE(time.Now()),
					job.Or(
						job.ReservedUntilIsNil(),
						job.ReservedUntilLT(time.Now()),
					),
				),
			).
			Order(job.ByCreatedAt()).
			ForUpdate(sql.WithLockAction(sql.SkipLocked)). // нет смысла ждать анлока записи, т.к. она уже выбрана
			First(ctx)
		if err != nil {
			return err
		}

		j, err = j.Update().
			SetAttempts(j.Attempts + 1).
			SetReservedUntil(until).
			Save(ctx)
		if err != nil {
			return err
		}
		retJob = Job{
			ID:       j.ID,
			Name:     j.Name,
			Payload:  j.Payload,
			Attempts: j.Attempts,
		}

		return nil
	})

	if err != nil {
		if store.IsNotFound(err) {
			return Job{}, ErrNoJobs
		}

		return Job{}, err
	}

	return retJob, nil
}

func (r *Repo) CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	j, err := r.db.Job(ctx).
		Create().
		SetName(name).
		SetPayload(payload).
		SetAvailableAt(availableAt).
		Save(ctx)
	if err != nil {
		return types.JobIDNil, err
	}

	return j.ID, nil
}

func (r *Repo) CreateFailedJob(ctx context.Context, name, payload, reason string) error {
	return r.db.FailedJob(ctx).
		Create().
		SetName(name).
		SetPayload(payload).
		SetReason(reason).
		Exec(ctx)
}

func (r *Repo) DeleteJob(ctx context.Context, jobID types.JobID) error {
	return r.db.Job(ctx).DeleteOneID(jobID).Exec(ctx)
}
