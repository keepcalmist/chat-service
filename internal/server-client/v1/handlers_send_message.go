package clientv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	internalErrors "github.com/keepcalmist/chat-service/internal/errors"
	"github.com/keepcalmist/chat-service/internal/middlewares"
	sendmessage "github.com/keepcalmist/chat-service/internal/usecases/client/send-message"
	"github.com/keepcalmist/chat-service/pkg/pointer"
)

func (h Handlers) PostSendMessage(eCtx echo.Context, params PostSendMessageParams) error {
	ctx := eCtx.Request().Context()

	reqBody := new(SendMessageRequest)
	err := eCtx.Bind(reqBody)
	if err != nil {
		return fmt.Errorf("bind request: %v", err)
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
		if errors.Is(err, sendmessage.ErrChatNotCreated) {
			return internalErrors.NewServerError(int(ErrorCodeCreateChatError), "h.sendMsg.Handle err", err)
		}
		if errors.Is(err, sendmessage.ErrProblemNotCreated) {
			return internalErrors.NewServerError(int(ErrorCodeCreateProblemError), "h.sendMsg.Handle err", err)
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
			AuthorId:  pointer.Ptr(resp.AuthorID),
			CreatedAt: resp.CreatedAt,
			Id:        resp.MessageID,
		},
	}
}
