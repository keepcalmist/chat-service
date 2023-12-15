package main

import (
	"context"

	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	"github.com/keepcalmist/chat-service/internal/types"
)

type dummyAdapter struct{}

func (dummyAdapter) Adapt(event eventstream.Event) (any, error) {
	return event, nil
}

type dummyEventStream struct{}

func (dummyEventStream) Subscribe(ctx context.Context, _ types.UserID) (<-chan eventstream.Event, error) {
	events := make(chan eventstream.Event)
	go func() {
		defer close(events)
		<-ctx.Done()
	}()
	return events, nil
}
