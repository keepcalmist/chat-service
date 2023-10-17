package clientv1

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/keepcalmist/chat-service/internal/types"
)

var stub = MessagesPage{Messages: []Message{
	{
		AuthorId:  types.NewUserID(),
		Body:      "Здравствуйте! Разберёмся.",
		CreatedAt: time.Now(),
		Id:        types.NewMessageID(),
	},
	{
		AuthorId:  types.MustParse[types.UserID]("f67e3424-ba2b-45ce-bf4e-e064f3663b78"),
		Body:      "Привет! Не могу снять денег с карты,\nпишет 'карта заблокирована'",
		CreatedAt: time.Now().Add(-time.Minute),
		Id:        types.NewMessageID(),
	},
}}

func (h Handlers) PostGetHistory(eCtx echo.Context, params PostGetHistoryParams) error {
	body, err := eCtx.Request().GetBody()
	if err != nil {
		return err
	}
	defer body.Close()

	return eCtx.JSON(http.StatusOK, GetHistoryResponse{
		Data: stub,
	})
}
