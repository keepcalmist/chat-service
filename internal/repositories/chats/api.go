package chatsrepo

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/keepcalmist/chat-service/internal/store/chat"
	"github.com/keepcalmist/chat-service/internal/types"
)

func (r *Repo) CreateIfNotExists(ctx context.Context, userID types.UserID) (types.ChatID, error) {
	err := r.db.Chat(ctx).
		Create().
		SetClientID(userID).
		OnConflict(sql.DoNothing()).Exec(ctx)
	if err != nil {
		return types.ChatIDNil, err
	}

	createdChat, err := r.db.Chat(ctx).
		Query().
		Where(
			chat.ClientID(userID),
		).First(ctx)
	if err != nil {
		return types.ChatIDNil, fmt.Errorf("get created chat: %w", err)
	}

	return createdChat.ID, nil
}
