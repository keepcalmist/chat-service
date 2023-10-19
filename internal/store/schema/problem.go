package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/keepcalmist/chat-service/internal/types"
)

// Problem holds the schema definition for the Problem entity.
type Problem struct {
	ent.Schema
}

// Fields of the Problem.
func (Problem) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.ProblemID{}).Immutable().Default(types.NewProblemID),
		field.UUID("manager_id", types.UserID{}).Immutable(),
		field.UUID("chat_id", types.ChatID{}).Immutable(),
		field.Time("created_at").Immutable().Default(time.Now().UTC()),
		field.Time("resolved_at").Nillable().Optional(),
	}
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("chat", Chat.Type).Required().Field("chat_id").Immutable().Unique(),
		edge.From("messages", Message.Type).Ref("problem"),
	}
}
