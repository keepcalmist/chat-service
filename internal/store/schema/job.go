package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/keepcalmist/chat-service/internal/types"
)

// jobMaxAttempts is some limit as protection from endless retries of outbox jobs.
const jobMaxAttempts = 30

type Job struct {
	ent.Schema
}

func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.JobID{}).Default(types.NewJobID).Unique().Immutable(),
		field.Text("name").NotEmpty().Immutable(),
		field.Text("payload").NotEmpty().Immutable(),
		field.Int("attempts").Default(0).Max(jobMaxAttempts),
		field.Time("available_at").Optional().Immutable(),
		field.Time("reserved_until").Optional(),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}

func (Job) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("created_at"),
		index.Fields("reserved_until", "available_at"),
	}
}

type FailedJob struct {
	ent.Schema
}

func (FailedJob) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.FailedJobID{}).Default(types.NewFailedJobID).Unique().Immutable(),
		field.Text("name").NotEmpty().Immutable(),
		field.Text("payload").NotEmpty().Immutable(),
		field.Text("reason").NotEmpty().Immutable(),
		field.Time("created_at").Immutable().Default(time.Now),
	}
}
