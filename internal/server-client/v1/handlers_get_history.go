package clientv1

import (
	"errors"
	"time"

	"github.com/keepcalmist/chat-service/internal/middlewares"
	gethistory "github.com/keepcalmist/chat-service/internal/usecases/client/get-history"
	"github.com/keepcalmist/chat-service/pkg/pointer"
	"github.com/labstack/echo/v4"

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
		// FIXME: 3) Обработать gethistory.ErrInvalidReqeest и gethistory.ErrInvalidCursor
		if errors.Is(err, gethistory.ErrInvalidRequest) {
			return
		}
	}

	// FIXME: 4) Сформировать ответ, обрабатывая возможное отсутствие автора у сообщения

	return nil
}
