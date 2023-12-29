package problemsrepo

import (
	"time"

	"github.com/keepcalmist/chat-service/internal/store"
	"github.com/keepcalmist/chat-service/internal/types"
)

type Problem struct {
	ID        types.ProblemID
	ChatID    types.ChatID
	CreatedAt time.Time
}

func adaptStoreProblems(problems []*store.Problem) []*Problem {
	adapted := make([]*Problem, 0, len(problems))
	for _, problem := range problems {
		adapted = append(adapted, adaptStoreProblem(problem))
	}

	return adapted
}

func adaptStoreProblem(p *store.Problem) *Problem {
	return &Problem{
		ID:        p.ID,
		ChatID:    p.ChatID,
		CreatedAt: p.CreatedAt,
	}
}
