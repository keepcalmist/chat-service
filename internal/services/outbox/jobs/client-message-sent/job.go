package clientmessagesentjob

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	"github.com/keepcalmist/chat-service/internal/services/outbox"
	"github.com/keepcalmist/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/job_mock.gen.go -package=clientmessagesentjobmocks

const Name = "client-message-sent"

type messageRepository interface {
	GetMessageByID(ctx context.Context, msgID types.MessageID) (*messagesrepo.Message, error)
	MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=job_options.gen.go -from-struct=Options
type Options struct {
	logger            *zap.Logger
	messageRepository messageRepository `option:"mandatory"  validate:"required"`
	eventStream       eventStream       `option:"mandatory" validate:"required"`
}

type Job struct {
	Options
	outbox.DefaultJob
}

func New(opts Options) (*Job, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("failed to create job <%s>: %w", Name, err)
	}

	if opts.logger == nil {
		opts.logger = zap.L().Named(Name)
	}

	return &Job{
		Options:    opts,
		DefaultJob: outbox.DefaultJob{},
	}, nil
}

func (j *Job) Name() string {
	return Name
}

func (j *Job) Handle(ctx context.Context, payload string) error {
	msgID := types.MessageID{}
	err := msgID.Scan(payload)
	if err != nil {
		j.logger.Error("failed to scan messageID", zap.Error(err))
		return fmt.Errorf("failed to scan payload in <%s> job: %w", Name, err)
	}

	j.logger.Info("handling blocked message", zap.String("message_id", msgID.String()))

	msg, err := j.messageRepository.GetMessageByID(ctx, msgID)
	if err != nil {
		return fmt.Errorf("failed to get message by id in <%s> job: %w", Name, err)
	}

	err = j.messageRepository.MarkAsVisibleForManager(ctx, msgID)
	if err != nil {
		return fmt.Errorf("failed to update block option of message in <%s> job:: %w", Name, err)
	}

	err = j.eventStream.Publish(ctx, msg.AuthorID, eventstream.NewMessageSentEvent(
		types.NewEventID(),
		msg.RequestID,
		msg.ID,
	))
	if err != nil {
		return fmt.Errorf("failed to publish event in <%s> job: %w", Name, err)
	}

	return nil
}
