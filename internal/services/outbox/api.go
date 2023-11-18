package outbox

import (
	"context"
	"time"

	"github.com/keepcalmist/chat-service/internal/types"
)

func (s *Service) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	return s.r.CreateJob(ctx, name, payload, availableAt)
}
