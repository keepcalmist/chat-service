package clientmessagesentjob_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	clientmessagesentjob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/client-message-sent"
	clientmessagesentjobmocks "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/client-message-sent/mocks"
	"github.com/keepcalmist/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgRepo := clientmessagesentjobmocks.NewMockmessageRepository(ctrl)
	streamer := clientmessagesentjobmocks.NewMockeventStream(ctrl)

	job, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(msgRepo, streamer))
	require.NoError(t, err)

	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	reqID := types.NewRequestID()
	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            clientID,
		RequestID:           reqID,
		Body:                body,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           false,
	}

	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)
	msgRepo.EXPECT().MarkAsVisibleForManager(gomock.Any(), msgID).Return(nil)

	streamer.EXPECT().Publish(gomock.Any(), clientID, eventMatcher{msg: msg}).Return(nil)

	// Action & assert.
	payload, err := MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}

var ErrInvalidMessageID = errors.New("invalid message id")

func MarshalPayload(messageID types.MessageID) (string, error) {
	if messageID == types.MessageIDNil {
		return "", ErrInvalidMessageID
	}

	return messageID.String(), nil
}

type eventMatcher struct {
	msg messagesrepo.Message
}

func (m eventMatcher) Matches(x interface{}) bool {
	if _, ok := x.(eventstream.Event); !ok {
		return false
	}

	switch t := x.(type) {
	case *eventstream.MessageSentEvent:
		if t.MessageID != m.msg.ID {
			return false
		}
		if t.RequestID != m.msg.RequestID {
			return false
		}
		if t.EventID == types.EventIDNil {
			return false
		}
		return true
	default:
		return false
	}
}

func (eventMatcher) String() string {
	return "eventstream.Event"
}
