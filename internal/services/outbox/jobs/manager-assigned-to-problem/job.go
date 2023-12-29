package managerassignedtoproblemjob

import (
	"context"
	"fmt"
	"time"

	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	msgproducer "github.com/keepcalmist/chat-service/internal/services/msg-producer"
	"github.com/keepcalmist/chat-service/internal/services/outbox"
	"github.com/keepcalmist/chat-service/internal/types"
	"go.uber.org/zap"
)

const Name = "manager-assigned-to-problem"

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

type problemsRepository interface {
	GetManagerID(ctx context.Context, problemID types.ProblemID) (types.UserID, error)
	GetClientId(ctx context.Context, problemID types.ProblemID) (types.UserID, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	msgRepo          problemsRepository `option:"mandatory"  validate:"required"`
	producer         messageProducer    `option:"mandatory"  validate:"required"`
	executionTimeout time.Duration      `option:"default=0"`
	maxAttempts      int                `option:"default=0"`
	logger           *zap.Logger
	eventStream      eventStream `option:"mandatory" validate:"required"`
}

type Job struct {
	Options
	defaultJob outbox.DefaultJob
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	if opts.logger == nil {
		opts.logger = zap.L().Named(Name)
	}

	return &Job{
		Options:    opts,
		defaultJob: outbox.DefaultJob{},
	}, nil
}

func (j *Job) Name() string {
	return Name
}

// FIXME: Реализуй джобу:
// FIXME: - клиенту улетает NewMessageEvent;
// FIXME: - менеджеру улетает NewChatEvent.

func (j *Job) Handle(ctx context.Context, payload string) error {
	problemId := types.ProblemID{}
	err := problemId.Scan(payload)
	if err != nil {
		return fmt.Errorf("failed to scan payload in <%s> job: %w", Name, err)
	}

	j.logger.Info("handling message", zap.String("message_id", msgID.String()))

	msg, err := j.msgRepo.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("failed to get message by id in <%s> job: %w", Name, err)
	}

	j.logger.Info("got message", zap.String("message_id", msgID.String()))

	msgToProduce := msgproducer.Message{
		ID:         msg.ID,
		ChatID:     msg.ChatID,
		Body:       msg.Body,
		FromClient: msg.IsVisibleForClient && !msg.IsService && !msg.IsBlocked && !msg.IsVisibleForManager,
	}

	err = j.producer.ProduceMessage(ctx, msgToProduce)
	if err != nil {
		return fmt.Errorf("failed to produce message in <%s> job: %w", Name, err)
	}

	j.logger.Info("message produced", zap.Stringer("message", msgToProduce))

	err = j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewNewMessageEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ChatID,
		msg.ID,
		msg.AuthorID,
		msg.CreatedAt,
		msg.Body,
		msg.IsService,
	))
	if err != nil {
		return fmt.Errorf("failed to publish event in <%s> job: %w", Name, err)
	}

	return nil
}

func (j *Job) ExecutionTimeout() time.Duration {
	if j.executionTimeout != time.Duration(0) {
		return j.executionTimeout
	}
	return j.defaultJob.ExecutionTimeout()
}

func (j *Job) MaxAttempts() int {
	if j.maxAttempts != 0 {
		return j.maxAttempts
	}
	return j.defaultJob.MaxAttempts()
}
