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

func adaptStoreProblems(ps []*store.Problem) []*Problem {
	problems := make([]*Problem, 0, len(ps))
	for _, p := range ps {
		problems = append(problems, adaptStoreProblem(p))
	}

	return problems
}

func adaptStoreProblem(p *store.Problem) *Problem {
	return &Problem{
		ID:        p.ID,
		ChatID:    p.ChatID,
		CreatedAt: p.CreatedAt,
	}
}
