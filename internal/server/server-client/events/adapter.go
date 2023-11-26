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
		err := adaptedEvent.FromMessageSentEvent(MessageSentEvent{
			EventId:   e.EventID,
			EventType: "MessageSentEvent",
			MessageId: e.MessageID,
			RequestId: e.RequestID,
		})
		if err != nil {
			return nil, err
		}
	case *eventstream.MessageBlockedEvent:
		err := adaptedEvent.FromMessageBlockedEvent(MessageBlockedEvent{
			EventId:   e.EventID,
			EventType: "MessageBlockedEvent",
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
			EventType: "MessageBlockedEvent",
			IsService: e.IsService,
			MessageId: e.MessageID,
			RequestId: e.RequestID,
		})
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("unknown event type: %T", ev)
	}

	return adaptedEvent, nil
}
