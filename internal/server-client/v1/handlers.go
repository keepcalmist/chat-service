package clientv1

import (
	"go.uber.org/zap"
)

//go:generate options-gen -out-filename=handlers_options.gen.go -from-struct=Options
type Options struct {
	logger *zap.Logger `option:"mandatory" validate:"required"`
	// Ждут своего часа.
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
