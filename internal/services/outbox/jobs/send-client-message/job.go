package sendclientmessagejob

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	msgproducer "github.com/keepcalmist/chat-service/internal/services/msg-producer"
	"github.com/keepcalmist/chat-service/internal/services/outbox"
	"github.com/keepcalmist/chat-service/internal/types"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=sendclientmessagejobmocks

const Name = "send-client-message"

type messageProducer interface {
	ProduceMessage(ctx context.Context, message msgproducer.Message) error
}

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	name             string            `option:"mandatory"  validate:"required"`
	msgRepo          messageRepository `option:"mandatory"  validate:"required"`
	producer         messageProducer   `option:"mandatory"  validate:"required"`
	executionTimeout time.Duration     `option:"default=0"`
	maxAttempts      int               `option:"default=0"`
	logger           *zap.Logger
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
		opts.logger = zap.L().Named(opts.name)
	}

	return &Job{
		Options:    opts,
		defaultJob: outbox.DefaultJob{},
	}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) error {
	msgID := types.MessageID{}
	err := json.Unmarshal([]byte(payload), &msgID)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload in <%s> job: %w", j.name, err)
	}

	j.logger.Info("handling message", zap.String("message_id", msgID.String()))

	msg, err := j.msgRepo.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("failed to get message by id in <%s> job: %w", j.name, err)
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
		return fmt.Errorf("failed to produce message in <%s> job: %w", j.Options.name, err)
	}

	j.logger.Info("message produced", zap.Stringer("message", msgToProduce))

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
