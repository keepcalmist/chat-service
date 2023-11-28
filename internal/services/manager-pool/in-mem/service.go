package inmemmanagerpool

import (
	"context"
	"sync"

	"go.uber.org/zap"

	managerpool "github.com/keepcalmist/chat-service/internal/services/manager-pool"
	"github.com/keepcalmist/chat-service/internal/types"
)

const (
	serviceName = "manager-pool"
	managersMax = 1000
)

type Service struct {
	logger   *zap.Logger
	ch       chan types.UserID
	contains map[types.UserID]struct{}
	mu       *sync.Mutex
}

func New() *Service {
	return &Service{
		logger:   zap.L().Named(serviceName),
		ch:       make(chan types.UserID, managersMax),
		contains: make(map[types.UserID]struct{}, managersMax),
		mu:       new(sync.Mutex),
	}
}

func (s *Service) Close() error {
	close(s.ch)
	return nil
}

func (s *Service) Get(_ context.Context) (types.UserID, error) {
	select {
	case managerID := <-s.ch:
		s.mu.Lock()
		delete(s.contains, managerID)
		s.mu.Unlock()
		return managerID, nil
	default:
		return types.UserIDNil, managerpool.ErrNoAvailableManagers
	}
}

func (s *Service) Put(_ context.Context, managerID types.UserID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.contains[managerID]; ok {
		return nil
	}

	s.contains[managerID] = struct{}{}
	s.ch <- managerID
	return nil
}

func (s *Service) Contains(_ context.Context, managerID types.UserID) (bool, error) {
	s.mu.Lock()
	_, ok := s.contains[managerID]
	s.mu.Unlock()

	return ok, nil
}

func (s *Service) Size() int {
	s.mu.Lock()
	length := len(s.contains)
	s.mu.Unlock()

	return length
}
