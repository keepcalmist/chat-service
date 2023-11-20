package problemsrepo

import (
	"context"
	"database/sql"
	"errors"

	entSql "entgo.io/ent/dialect/sql"

	"github.com/keepcalmist/chat-service/internal/store/problem"
	"github.com/keepcalmist/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, chatID types.ChatID) (types.ProblemID, error) {
	id, err := r.db.Problem(ctx).
		Create().
		SetChatID(chatID).
		OnConflict(entSql.DoNothing()).DoNothing().ID(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return types.ProblemIDNil, err
	}

	if id != types.ProblemIDNil {
		return id, nil
	}

	createdProblem, err := r.db.Problem(ctx).
		Query().
		Unique(false).
		Where(
			problem.ChatID(chatID),
			problem.ResolvedAtIsNil(),
		).Order(problem.ByCreatedAt(func(options *entSql.OrderTermOptions) {
		options.Desc = true
	})).First(ctx)
	if err != nil {
		return types.ProblemIDNil, err
	}

	return createdProblem.ID, nil
}
