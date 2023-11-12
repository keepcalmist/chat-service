package clientv1

import (
	"errors"
	"net/http"

	internalErrors "github.com/keepcalmist/chat-service/internal/errors"
	"github.com/keepcalmist/chat-service/internal/middlewares"
	sendmessage "github.com/keepcalmist/chat-service/internal/usecases/client/send-message"
	"github.com/keepcalmist/chat-service/pkg/pointer"
	"github.com/labstack/echo/v4"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()

	reqBody := new(SendMessageRequest)
	err := eCtx.Bind(reqBody)
	if err != nil {
		return err
	}

	clientID, ok := middlewares.GetUserID(eCtx)
	if !ok {
		return internalErrors.NewServerError(http.StatusBadRequest, "cannot get clientID from context", nil)
	}

	resp, err := h.sendMsg.Handle(ctx, sendmessage.Request{
		ID:          params.XRequestID,
		ClientID:    clientID,
		MessageBody: reqBody.MessageBody,
	})

	if err != nil {
		if errors.Is(err, sendmessage.ErrInvalidRequest) {
			return internalErrors.NewServerError(http.StatusBadRequest, "h.sendMsg.Handle err", err)
		}

		return internalErrors.NewServerError(http.StatusInternalServerError, "h.sendMsg.Handle err", err)
	}

	err = eCtx.JSONPretty(http.StatusOK, convertSendMessageResponseToMessage(resp), "  ")
	if err != nil {
		return err
	}

	return nil
}

func convertSendMessageResponseToMessage(resp sendmessage.Response) SendMessageResponse {
	return SendMessageResponse{
		Data: &MessageHeader{
			AuthorID:  pointer.Ptr(resp.AuthorID),
			CreatedAt: resp.CreatedAt,
			Id:        resp.MessageID,
		},
	}
}
