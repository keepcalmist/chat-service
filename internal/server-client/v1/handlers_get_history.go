package clientv1

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	internalErrors "github.com/keepcalmist/chat-service/internal/errors"
	"github.com/keepcalmist/chat-service/internal/middlewares"
	gethistory "github.com/keepcalmist/chat-service/internal/usecases/client/get-history"
	"github.com/keepcalmist/chat-service/pkg/pointer"
)

func (h Handlers) PostGetHistory(eCtx echo.Context, req PostGetHistoryParams) error {
	ctx := eCtx.Request().Context()

	reqBody := new(GetHistoryRequest)
	err := eCtx.Bind(reqBody)
	if err != nil {
		return fmt.Errorf("bind request: %w", err)
	}

	clientID, ok := middlewares.GetUserID(eCtx)
	if !ok {
		return internalErrors.NewServerError(http.StatusBadRequest, "cannot get clientID from context", nil)
	}

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
		return internalErrors.NewServerError(http.StatusInternalServerError, "h.getHistory.Handle err", err)
	}

	err = eCtx.JSONPretty(http.StatusOK, adaptGetHistoryResponse(resp), "  ")
	if err != nil {
		return fmt.Errorf("JSONPretty: %w", err)
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
		AuthorId:   pointer.PtrWithZeroAsNil(message.AuthorID),
		Body:       message.Body,
		CreatedAt:  message.CreatedAt,
		Id:         message.ID,
		IsBlocked:  message.IsBlocked,
		IsReceived: message.IsReceived,
		IsService:  message.IsService,
	}
}
