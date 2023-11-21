package sendclientmessagejob

import (
	"errors"

	"github.com/keepcalmist/chat-service/internal/types"
)

var ErrInvalidMessageID = errors.New("invalid message id")

func MarshalPayload(messageID types.MessageID) (string, error) {
	if messageID == types.MessageIDNil {
		return "", ErrInvalidMessageID
	}

	return messageID.String(), nil
}
