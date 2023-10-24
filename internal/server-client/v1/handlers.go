package clientv1

import (
	"context"

	gethistory "github.com/keepcalmist/chat-service/internal/usecases/client/get-history"
	"go.uber.org/zap"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/handlers_mocks.gen.go -package=clientv1mocks

type getHistoryUseCase interface {
	Handle(ctx context.Context, req gethistory.Request) (gethistory.Response, error)
}

//go:generate options-gen -out-filename=handlers_options.gen.go -from-struct=Options
type Options struct {
	getHistory getHistoryUseCase `option:"mandatory" validate:"required"`
	logger     *zap.Logger       `option:"mandatory" validate:"required"`
}

type Handlers struct {
	Options
}

func NewHandlers(opts Options) (Handlers, error) {
	if err := opts.Validate(); err != nil {
		return Handlers{}, err
	}
	return Handlers{Options: opts}, nil
}
