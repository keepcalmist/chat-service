package chatsrepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	entSql "entgo.io/ent/dialect/sql"

	"github.com/keepcalmist/chat-service/internal/store/chat"
	"github.com/keepcalmist/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	err := r.db.Chat(ctx).
		Create().
		SetClientID(userID).
		OnConflict(entSql.DoNothing()).
		Exec(ctx)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return types.ChatIDNil, err
	}

	createdChat, err := r.db.Chat(ctx).
		Query().
		Where(
			chat.ClientID(userID),
		).Unique(false).First(ctx)
	if err != nil {
		return types.ChatIDNil, fmt.Errorf("get created chat: %w", err)
	}

	return createdChat.ID, nil
}
