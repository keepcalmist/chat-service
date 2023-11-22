package errhandler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	internalErrors "github.com/keepcalmist/chat-service/internal/errors"
)

var _ echo.HTTPErrorHandler = Handler{}.Handle

//go:generate options-gen -out-filename=errhandler_options.gen.go -from-struct=Options
type Options struct {
	logger          *zap.Logger                                    `option:"mandatory" validate:"required"`
	productionMode  bool                                           `option:"mandatory"`
	responseBuilder func(code int, msg string, details string) any `option:"mandatory" validate:"required"`
}

type Handler struct {
	lg              *zap.Logger
	productionMode  bool
	responseBuilder func(code int, msg string, details string) any
}

func New(opts Options) (Handler, error) {
	if err := opts.Validate(); err != nil {
		return Handler{}, err
	}

	return Handler{
		lg:              opts.logger,
		productionMode:  opts.productionMode,
		responseBuilder: opts.responseBuilder,
	}, nil
}

func (h Handler) Handle(err error, eCtx echo.Context) {
	code, msg, details := internalErrors.ProcessServerError(err)

	h.lg.Error("server error", []zap.Field{
		zap.Int("code", code),
		zap.String("msg", msg),
		zap.String("details", details),
	}...)
	if h.productionMode {
		details = ""
	}

	err = eCtx.JSONPretty(http.StatusOK, h.responseBuilder(code, msg, details), "  ")
	if err != nil {
		h.lg.Error("cannot write response", zap.Error(err))
		eCtx.Error(err)
	}
}
