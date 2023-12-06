package messagesrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/keepcalmist/chat-service/internal/store/message"
	"github.com/keepcalmist/chat-service/internal/types"
)

func (r *Repo) MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error {
	qb := r.db.Message(ctx).
		Update().
		SetIsVisibleForManager(true).
		SetCheckedAt(time.Now()).
		Where(message.ID(msgID))

	err := qb.Exec(ctx)
	if err != nil {
		return fmt.Errorf("Repo.MarkAsVisibleForManager err: %w", err)
	}

	return nil
}

func (r *Repo) BlockMessage(ctx context.Context, msgID types.MessageID) error {
	qb := r.db.Message(ctx).
		Update().
		SetIsBlocked(true).
		SetCheckedAt(time.Now()).
		Where(message.ID(msgID))

	err := qb.Exec(ctx)
	if err != nil {
		return fmt.Errorf("Repo.BlockMessage err: %w", err)
	}

	return nil
}
