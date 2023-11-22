package canreceiveproblems

import (
	"context"
	"fmt"

	"github.com/keepcalmist/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=canreceiveproblemsmocks

type managerLoadService interface {
	CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error)
}

type managerPool interface {
	Contains(ctx context.Context, managerID types.UserID) (bool, error)
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

func (u UseCase) Handle(ctx context.Context, req Request) (Response, error) {
	if err := req.Validate(); err != nil {
		return Response{}, fmt.Errorf("validation request err: %w", err)
	}

	contains, err := u.managerPool.Contains(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("checking manager %d in pool err: %w", req.ManagerID, err)
	}

	if contains {
		return Response{Result: false}, nil
	}

	canTakeProblem, err := u.managerLoadService.CanManagerTakeProblem(ctx, req.ManagerID)
	if err != nil {
		return Response{}, fmt.Errorf("can manager %d take problem err: %w", req.ManagerID, err)
	}

	return Response{
		Result: canTakeProblem,
	}, nil
}
