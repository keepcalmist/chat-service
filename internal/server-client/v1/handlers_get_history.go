package clientv1

import (
	"errors"
	"net/http"
	"time"

	"github.com/keepcalmist/chat-service/internal/middlewares"
	gethistory "github.com/keepcalmist/chat-service/internal/usecases/client/get-history"
	"github.com/keepcalmist/chat-service/pkg/pointer"
	"github.com/labstack/echo/v4"

	internalErrors "github.com/keepcalmist/chat-service/internal/errors"
	"github.com/keepcalmist/chat-service/internal/types"
)

var stub = MessagesPage{Messages: []Message{
	{
		AuthorId:  pointer.Ptr(types.NewUserID()),
		Body:      "Здравствуйте! Разберёмся.",
		CreatedAt: time.Now(),
		Id:        types.NewMessageID(),
	},
	{
		AuthorId:  pointer.Ptr(types.MustParse[types.UserID]("f67e3424-ba2b-45ce-bf4e-e064f3663b78")),
		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
		CreatedAt: time.Now().Add(-time.Minute),
		Id:        types.NewMessageID(),
	},
}}

func (h Handlers) PostGetHistory(eCtx echo.Context, req PostGetHistoryParams) error {
	ctx := eCtx.Request().Context()
	clientID := middlewares.MustUserID(eCtx)

	// FIXME: 1) За-bind-ить входящий запрос
	reqBody := new(GetHistoryRequest)
	err := eCtx.Bind(reqBody)
	if err != nil {
		return err
	}

	// FIXME: 2) Вызвать соответствующий юзкейс
	resp, err := h.getHistory.Handle(ctx, gethistory.Request{
		ID:       req.XRequestID,
		ClientID: clientID,
		PageSize: pointer.Indirect(reqBody.PageSize),
		Cursor:   pointer.Indirect(reqBody.Cursor),
	})
	if err != nil {
		if errors.Is(err, gethistory.ErrInvalidRequest) || errors.Is(err, gethistory.ErrInvalidCursor) {
			return internalErrors.NewServerError(http.StatusBadRequest, "h.getHistory.Handle err", err)
		}
	}

	err = eCtx.JSONPretty(http.StatusOK, adaptGetHistoryResponse(resp), "  ")
	if err != nil {
		return err
	}

	return nil
}

func adaptGetHistoryResponse(resp gethistory.Response) GetHistoryResponse {
	return GetHistoryResponse{
		Data: &MessagesPage{
			Messages: adaptMessages(resp.Messages),
			Next:     resp.NextCursor,
		},
	}
}

func adaptMessages(messages []gethistory.Message) []Message {
	arr := make([]Message, 0, len(messages))
	for _, msg := range messages {
		arr = append(arr, adaptMessage(msg))
	}
	return arr
}

func adaptMessage(message gethistory.Message) Message {
	return Message{
		AuthorId:  pointer.PtrWithZeroAsNil(message.AuthorID),
		Body:      message.Body,
		CreatedAt: message.CreatedAt,
		Id:        message.ID,
	}
}
