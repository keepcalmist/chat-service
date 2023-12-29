package managerevents

import (
	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	websocketstream "github.com/keepcalmist/chat-service/internal/websocket-stream"
)

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	// FIXME: Реализуй меня.
	// FIXME: Покрой юнит-тестами.
	return nil, nil
}
