package gethistory

import (
	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
)

func convertMessages(m []messagesrepo.Message) []Message {
	convertedMessages := make([]Message, 0, len(m))
	for _, msg := range m {
		convertedMessages = append(convertedMessages, convertMessage(msg))
	}

	return convertedMessages
}

func convertMessage(m messagesrepo.Message) Message {
	return Message{
		ID:         m.ID,
		AuthorID:   m.AuthorID,
		Body:       m.Body,
		CreatedAt:  m.CreatedAt,
		IsBlocked:  m.IsBlocked,
		IsService:  m.IsService,
		IsReceived: m.IsVisibleForManager && !m.IsBlocked,
	}
}
