package messagesrepo

import (
	"time"

	"github.com/keepcalmist/chat-service/internal/store"
	"github.com/keepcalmist/chat-service/internal/types"
)

type Message struct {
	ID                  types.MessageID
	ChatID              types.ChatID
	AuthorID            types.UserID
	RequestID           types.RequestID
	Body                string
	CreatedAt           time.Time
	IsVisibleForClient  bool
	IsVisibleForManager bool
	IsBlocked           bool
	IsService           bool
}

func adaptStoreMessage(m *store.Message) Message {
	return Message{
		ID:                  m.ID,
		ChatID:              m.ChatID,
		AuthorID:            m.AuthorID,
		RequestID:           m.InitialRequestID,
		Body:                m.Body,
		CreatedAt:           m.CreatedAt,
		IsVisibleForClient:  m.IsVisibleForClient,
		IsVisibleForManager: m.IsVisibleForManager,
		IsBlocked:           m.IsBlocked,
		IsService:           m.IsService,
	}
}
