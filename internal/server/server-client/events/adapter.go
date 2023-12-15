package clientevents

import (
	"fmt"

	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	websocketstream "github.com/keepcalmist/chat-service/internal/websocket-stream"
	"github.com/keepcalmist/chat-service/pkg/pointer"
)

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	if err := ev.Validate(); err != nil {
		return nil, err
	}

	adaptedEvent := new(Event)

	switch e := ev.(type) {
	case *eventstream.MessageSentEvent:

		err := adaptedEvent.FromMessageSentEvent(MessageId{
			EventId:   e.EventID,
			MessageId: e.MessageID,
			RequestId: e.RequestID,
		})
		if err != nil {
			return nil, err
		}

	case *eventstream.MessageBlockedEvent:
		err := adaptedEvent.FromMessageBlockedEvent(MessageId{
			EventId:   e.EventID,
			MessageId: e.MessageID,
			RequestId: e.RequestID,
		})
		if err != nil {
			return nil, err
		}

	case *eventstream.NewMessageEvent:
		err := adaptedEvent.FromNewMessageEvent(NewMessageEvent{
			AuthorId:  pointer.PtrWithZeroAsNil(e.AuthorID),
			Body:      e.Body,
			CreatedAt: e.CreatedAt,
			EventId:   e.EventID,
			IsService: e.IsService,
			MessageId: e.MessageID,
			RequestId: e.RequestID,
		})
		if err != nil {
			return nil, err
		}
	case *eventstream.MessageID:
		err := adaptedEvent.FromMessageId(MessageId{
			MessageId: e.MessageID,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown event type: %T", ev)
	}

	return adaptedEvent, nil
}
