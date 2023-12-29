package problemsrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	entSql "entgo.io/ent/dialect/sql"
	"github.com/keepcalmist/chat-service/internal/store"
	"github.com/keepcalmist/chat-service/internal/store/chat"
	"github.com/keepcalmist/chat-service/internal/store/message"

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

func (r *Repo) GetManagerOpenProblemsCount(ctx context.Context, managerID types.UserID) (int, error) {
	return r.db.Problem(ctx).
		Query().
		Unique(false).
		Where(
			problem.ManagerID(managerID),
			problem.ResolvedAtIsNil(),
		).Count(ctx)
}

func (r *Repo) GetUnassignedProblems(ctx context.Context) ([]*Problem, error) {
	problems, err := r.db.Problem(ctx).
		Query().
		Where(
			problem.ManagerIDIsNil(),
			problem.ResolvedAtIsNil(),
		).Order(store.Asc(problem.FieldCreatedAt)).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unassigned problems: %w", err)
	}

	return adaptStoreProblems(problems), nil
}

func (r *Repo) SetManagerForProblem(
	ctx context.Context,
	problemID types.ProblemID,
	managerID types.UserID,
) error {
	err := r.db.Problem(ctx).UpdateOneID(problemID).SetManagerID(managerID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to set manager for problem: %w", err)
	}

	return nil
}

func (r *Repo) GetManagerID(ctx context.Context, problemID types.ProblemID) (types.UserID, error) {
	p, err := r.db.Problem(ctx).Get(ctx, problemID)
	if err != nil {
		return types.UserIDNil, fmt.Errorf("failed to set manager for problem: %w", err)
	}

	return *p.ManagerID, nil
}

func (r *Repo) GetClientId(ctx context.Context, problemID types.ProblemID) (types.UserID, error) {
	chatWithProblem, err := r.db.Chat(ctx).
		Query().
		Where(chat.HasProblemsWith(problem.IDEQ(problemID))).
		Only(ctx)
	if err != nil {
		return types.UserIDNil, fmt.Errorf("failed to set manager for problem: %w", err)
	}

	return chatWithProblem.ClientID, nil
}

func (r *Repo) GetRequestID(
	ctx context.Context,
	problemID types.ProblemID,
) (types.RequestID, error) {
	p, err := r.db.Problem(ctx).Query().Where(problem.ID(problemID)).Only(ctx)
	if err != nil {
		return types.RequestIDNil, fmt.Errorf("failed to get problem: %w", err)
	}
	chat, err := r.db.Chat(ctx).Get(ctx, p.ChatID)
	if err != nil {
		return types.RequestIDNil, fmt.Errorf("failed to get chat: %w", err)
	}

	msg, err := r.db.Message(ctx).
		Query().
		Unique(true).
		Where(message.ChatIDEQ(chat.ID)).
		Order(store.Desc(message.FieldID)).
		First(ctx)
	if err != nil {
		return types.RequestIDNil, fmt.Errorf("failed to get messages: %w", err)
	}

	return msg.InitialRequestID, nil
}
