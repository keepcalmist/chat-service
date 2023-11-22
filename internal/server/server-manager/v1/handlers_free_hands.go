package managerv1

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	internalErrors "github.com/keepcalmist/chat-service/internal/errors"
	"github.com/keepcalmist/chat-service/internal/middlewares"
	canreceiveproblems "github.com/keepcalmist/chat-service/internal/usecases/manager/can-receive-problems"
)

func (h Handlers) PostGetFreeHandsBtnAvailability(eCtx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error {
	ctx := eCtx.Request().Context()

	managerID, ok := middlewares.GetUserID(eCtx)
	if !ok {
		return internalErrors.NewServerError(http.StatusBadRequest, "cannot get managerID from context", nil)
	}

	resp, err := h.canReceiveProblemsUseCase.Handle(ctx, canreceiveproblems.Request{
		ID:        params.XRequestID,
		ManagerID: managerID,
	})
	if err != nil {
		return internalErrors.NewServerError(http.StatusInternalServerError, "h.canReceiveProblemsUseCase.Handle err", err)
	}

	err = eCtx.JSONPretty(http.StatusOK, adaptGetFreeHandsBtnAvailabilityResponse(resp), "  ")
	if err != nil {
		return internalErrors.NewServerError(http.StatusInternalServerError, fmt.Sprintf("JSONPretty: %v", err), nil)
	}

	return nil
}

func adaptGetFreeHandsBtnAvailabilityResponse(resp canreceiveproblems.Response) GetFreeHandsBtnAvailabilityResponse {
	return GetFreeHandsBtnAvailabilityResponse{
		Data: &GetFreeHandsBtnAvailability{
			Available: resp.Result,
		},
	}
}
