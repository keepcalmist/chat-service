package inmemeventstream

import (
	"context"
	"errors"
	"sync/atomic"
	"time"

	"github.com/imkira/go-observer/v2"
	"go.uber.org/zap"

	eventstream "github.com/keepcalmist/chat-service/internal/services/event-stream"
	"github.com/keepcalmist/chat-service/internal/types"
)

const (
	serviceName = "event-stream"
	timeToClose = 30 * time.Second
	timeToWrite = 10 * time.Second
)

type Service struct {
	prop map[types.UserID]observer.Property[eventstream.Event]

	logger *zap.Logger
	count  atomic.Int32
}

func New() *Service {
	return &Service{
		prop:   make(map[types.UserID]observer.Property[eventstream.Event]),
		logger: zap.L().Named(serviceName),
	}
}

func (s *Service) Subscribe(ctx context.Context, userID types.UserID) (<-chan eventstream.Event, error) {
	eventsChan := make(chan eventstream.Event, 10)
	prop := observer.Property[eventstream.Event](nil)
	if _, ok := s.prop[userID]; !ok {
		s.prop[userID] = observer.NewProperty[eventstream.Event](nil)
	}

	prop = s.prop[userID]

	stream := prop.Observe()
	go func() {
		s.count.Add(1)
		defer s.count.Add(-1)
		defer close(eventsChan)

		for {
			select {
			case <-ctx.Done():
				return
			case <-stream.Changes():
				stream.Next()
			}

			event := stream.Value()
			if event == nil {
				s.logger.Error("event is nil")
				select {
				case <-ctx.Done():
					return
				case <-stream.Changes():
					stream.Next()
				}
				continue
			}

			if err := event.Validate(); err != nil {
				s.logger.Error("event validation failed", zap.Error(event.Validate()))
				continue
			}

			select {
			case <-time.After(timeToWrite):
				s.logger.Error("writing to chan timeout exceeded", zap.Duration("timeout", timeToWrite))
				close(eventsChan)
				return
			case eventsChan <- event:
			}
		}
	}()

	return eventsChan, nil
}

func (s *Service) Publish(_ context.Context, userID types.UserID, event eventstream.Event) error {
	err := event.Validate()
	if err != nil {
		return err
	}

	if _, ok := s.prop[userID]; !ok {
		return nil
	}

	s.prop[userID].Update(event)

	return nil
}

func (s *Service) Close() error {
	select {
	case <-time.After(timeToClose):
		s.logger.Error("timeout exceeded", zap.Duration("timeout", timeToClose))
		return errors.New("timeout exceeded")
	default:
		if s.count.Load() == 0 {
			return nil
		}
	}

	return nil
}
