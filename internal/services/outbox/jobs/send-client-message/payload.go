package sendclientmessagejob

import (
	"github.com/keepcalmist/chat-service/internal/types"
)

// FIXME: Вероятно необходимо добавить приватных типов и функций.

func MarshalPayload(messageID types.MessageID) (string, error) {
	str, err := messageID.MarshalText()
	if err != nil {
		return "", err
	}

	return string(str), nil
}
