package managerload

import (
	"context"

	"github.com/keepcalmist/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=managerloadmocks

type problemsRepository interface {
	GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error)
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	maxProblemsAtTime int `option:"mandatory" validate:"required,min=1,max=30"`

	problemsRepo problemsRepository `option:"mandatory" validate:"required"`
}

type Service struct {
	Options
}

func New(opts Options) (*Service, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	return &Service{Options: opts}, nil
}
