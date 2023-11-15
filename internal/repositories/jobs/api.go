package jobsrepo

import (
	"context"
	"errors"
	"time"

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
	// FIXME: Избегая гонки на уровне строчек БД, сделать следующее:
	// FIXME: - выбрать не зарезервированную другим воркером джобу, чьё время выполнения уже настало;
	// FIXME: - увеличить счётчик попыток выполнения на 1;
	// FIXME: - зарезервировать джобу до `until`;
	// FIXME: - вернуть в ответ необходимую инфу о джобе.
	r.db.Job(ctx).Query().Select(job.FieldID).Where().Select
	reservedJob, err :=
	if err != nil {
		return Job{}, err
	}

	return Job{}, nil
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
