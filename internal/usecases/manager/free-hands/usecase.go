package freehands

import (
	"context"
	"errors"
	"fmt"

	"github.com/keepcalmist/chat-service/internal/types"
)

var ErrManagerCannotTakeMoreProblems = errors.New("manager cannot take more problems")

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=freehandsmocks

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Put(ctx context.Context, managerID types.UserID) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	managerLoadService managerLoadService `option:"mandatory" validate:"required"`
	managerPool        managerPool        `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, err
	}
	return UseCase{
		Options: opts,
	}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validation request err: %w", err)
	}

	canTakeProblem, err := u.managerLoadService.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return fmt.Errorf("can manager %d take problem err: %w", req.ManagerID, err)
	}

	if !canTakeProblem {
		return ErrManagerCannotTakeMoreProblems
	}

	return u.managerPool.Put(ctx, req.ManagerID)
}
