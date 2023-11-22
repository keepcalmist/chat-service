package problemsrepo

import (
	"context"

	"entgo.io/ent/dialect/sql"

	"github.com/keepcalmist/chat-service/internal/store/problem"
	"github.com/keepcalmist/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	err := r.db.Problem(ctx).
		Create().
		SetChatID(chatID).
		OnConflict(sql.DoNothing()).
		Exec(ctx)
	if err != nil {
		return types.ProblemIDNil, err
	}

	createdProblem, err := r.db.Problem(ctx).
		Query().
		Unique(false).
		Where(
			problem.ChatID(chatID),
		).First(ctx)
	if err != nil {
		return types.ProblemIDNil, err
	}

	return createdProblem.ID, nil
}
