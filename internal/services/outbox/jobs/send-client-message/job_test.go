package sendclientmessagejob_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	msgproducer "github.com/keepcalmist/chat-service/internal/services/msg-producer"
	sendclientmessagejob "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/send-client-message"
	sendclientmessagejobmocks "github.com/keepcalmist/chat-service/internal/services/outbox/jobs/send-client-message/mocks"
	"github.com/keepcalmist/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendclientmessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendclientmessagejobmocks.NewMockmessageRepository(ctrl)
	eventStream := sendclientmessagejobmocks.NewMockeventStream(ctrl)
	job, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(msgRepo, msgProducer, eventStream))
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

	msgProducer.EXPECT().ProduceMessage(gomock.Any(), msgproducer.Message{
		ID:         msgID,
		ChatID:     chatID,
		Body:       body,
		FromClient: true,
	}).Return(nil)

	eventStream.EXPECT().Publish(gomock.Any(), clientID, eventMatcher{msg: msg}).Return(nil)

	// Action & assert.
	payload, err := sendclientmessagejob.MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}

type eventMatcher struct {
	msg messagesrepo.Message
}

func (m eventMatcher) Matches(x interface{}) bool {
	if _, ok := x.(eventstream.Event); !ok {
		return false
	}

	switch t := x.(type) {
	case *eventstream.NewMessageEvent:
		if t.ChatID != m.msg.ChatID {
			return false
		}
		if t.MessageID != m.msg.ID {
			return false
		}
		if t.RequestID != m.msg.RequestID {
			return false
		}
		if t.Body != m.msg.Body {
			return false
		}
		if t.AuthorID != m.msg.AuthorID {
			return false
		}
		if t.IsService != m.msg.IsService {
			return false
		}
		if t.CreatedAt != m.msg.CreatedAt {
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
