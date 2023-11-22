package managerv1

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	internalErrors "github.com/keepcalmist/chat-service/internal/errors"
	"github.com/keepcalmist/chat-service/internal/middlewares"
	freehands "github.com/keepcalmist/chat-service/internal/usecases/manager/free-hands"
	"github.com/keepcalmist/chat-service/pkg/pointer"
)

func (h Handlers) PostFreeHands(eCtx echo.Context, params PostFreeHandsParams) error {
	ctx := eCtx.Request().Context()

	managerID, ok := middlewares.GetUserID(eCtx)
	if !ok {
		return internalErrors.NewServerError(http.StatusBadRequest, "cannot get managerID from context", nil)
	}

	err := h.freeHandsUseCase.Handle(ctx, freehands.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		if errors.Is(err, freehands.ErrManagerCannotTakeMoreProblems) {
			return internalErrors.NewServerError(int(ErrorManagerCannotTakeMoreProblems), "manager cannot take more problems", err)
		}
		return internalErrors.NewServerError(http.StatusInternalServerError, "h.freeHandsUseCase.Handle err", err)
	}

	err = eCtx.JSONPretty(http.StatusOK, FreeHandsResponse{
		Data: pointer.Ptr(make(map[string]interface{})),
	}, "  ")
	if err != nil {
		return internalErrors.NewServerError(http.StatusInternalServerError, "JSONPretty err", err)
	}

	return nil
}
