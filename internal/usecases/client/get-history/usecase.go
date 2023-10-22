package gethistory

import (
	"context"
	"errors"

	"github.com/keepcalmist/chat-service/internal/cursor"
	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	"github.com/keepcalmist/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=gethistorymocks

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrInvalidCursor  = errors.New("invalid cursor")
)

type messagesRepository interface {
	GetClientChatMessages(
		ctx context.Context,
		clientID types.UserID,
		pageSize int,
		cursor *messagesrepo.Cursor,
	) ([]messagesrepo.Message, *messagesrepo.Cursor, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	msgRepo messagesRepository `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, err
	}

	return UseCase{
		Options: opts,
	}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if req.Validate() != nil {
		return Response{}, ErrInvalidRequest
	}

	var cur *messagesrepo.Cursor
	if req.Cursor != "" {
		cur = new(messagesrepo.Cursor)
		err := cursor.Decode(req.Cursor, cur)
		if err != nil {
			return Response{}, ErrInvalidCursor
		}
	}

	msgs, next, err := u.msgRepo.GetClientChatMessages(ctx, req.ClientID, req.PageSize, cur)
	if err != nil {
		if errors.Is(err, messagesrepo.ErrInvalidCursor) {
			return Response{}, ErrInvalidCursor
		}
		return Response{}, err
	}

	nextCursor := ""
	if next != nil {
		nextCursor, err = cursor.Encode(next)
		if err != nil {
			return Response{}, err
		}
	}

	return Response{
		Messages:   convertMessages(msgs),
		NextCursor: nextCursor,
	}, nil
}
