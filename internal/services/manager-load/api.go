package managerload

import (
	"context"
	"fmt"

	"github.com/keepcalmist/chat-service/internal/types"
)

func (s *Service) CanManagerTakeProblem(ctx context.Context, managerID types.UserID) (bool, error) {
	count, err := s.problemsRepo.GetManagerOpenProblemsCount(ctx, managerID)
	if err != nil {
		return false, fmt.Errorf("failed to get counf of manager's open problems: %w", err)
	}

	return s.maxProblemsAtTime > count, nil
}
