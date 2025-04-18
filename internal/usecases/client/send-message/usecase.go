package sendmessage

import (
	"context"
	"errors"
	"fmt"
	"time"

	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	sendclientmessagejob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/send-client-message"
	"github.com/keepcalmist/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=sendmessagemocks

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrChatNotCreated    = errors.New("chat not created")
	ErrProblemNotCreated = errors.New("problem not created")
)

type chatsRepository interface {
	CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error)
}

type messagesRepository interface {
	GetMessageByRequestID(ctx context.Context, reqID types.RequestID) (*messagesrepo.Message, error)
	CreateClientVisible(
		ctx context.Context,
		reqID types.RequestID,
		problemID types.ProblemID,
		chatID types.ChatID,
		authorID types.UserID,
		msgBody string,
	) (*messagesrepo.Message, error)
}

type problemsRepository interface {
	CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error)
}

type keepcalmist/chat-servicetor interface {
	RunInTx(ctx context.Context, f func(context.Context) error) error
}

type outboxService interface {
	Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	chatRepo    chatsRepository    `option:"mandatory" validate:"required"`
	msgRepo     messagesRepository `option:"mandatory" validate:"required"`
	outbox      outboxService      `option:"mandatory" validate:"required"`
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	tx          transactor         `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	return UseCase{Options: opts}, opts.Validate()
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if req.Validate() != nil {
		return Response{}, ErrInvalidRequest
	}
	resp := Response{}

	err := u.tx.RunInTx(ctx, func(ctx context.Context) error {
		msg, err := u.msgRepo.GetMessageByRequestID(ctx, req.ID)
		if err != nil && err != messagesrepo.ErrMsgNotFound {
			return err
		}

		if msg != nil {
			resp = convertResponse(msg)
			return nil
		}

		chatID, err := u.chatRepo.CreateIfNotExists(ctx, req.ClientID)
		if err != nil {
			return ErrChatNotCreated
		}

		problemID, err := u.problemRepo.CreateIfNotExists(ctx, chatID)
		if err != nil {
			return ErrProblemNotCreated
		}

		msg, err = u.msgRepo.CreateClientVisible(
			ctx,
			req.ID,
			problemID,
			chatID,
			req.ClientID,
			req.MessageBody,
		)
		if err != nil {
			return err
		}

		_, err = u.outbox.Put(ctx, sendclientmessagejob.Name, msg.ID.String(), time.Now())
		if err != nil {
			return fmt.Errorf("failed to put job to outbox: %w", err)
		}

		resp = convertResponse(msg)

		return nil
	})
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}

func convertResponse(msg *messagesrepo.Message) Response {
	return Response{
		AuthorID:  msg.AuthorID,
		MessageID: msg.ID,
		CreatedAt: msg.CreatedAt,
	}
}
