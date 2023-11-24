package shutdown

import "sync"

type ShutDown struct {
	once *sync.Once
	ch   chan struct{}
}

func NewShutDown() *ShutDown {
	return &ShutDown{
		once: &sync.Once{},
		ch:   make(chan struct{}),
	}
}

func (s *ShutDown) Close() error {
	s.once.Do(func() {
		close(s.ch)
	})

	return nil
}

func (s *ShutDown) Done() <-chan struct{} {
	return s.ch
}
