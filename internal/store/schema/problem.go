package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

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
		field.UUID("manager_id", types.UserID{}).Nillable().Optional(),
		field.UUID("chat_id", types.ChatID{}).Immutable(),
		field.Time("created_at").Immutable().Default(time.Now),
		field.Time("resolved_at").Nillable().Optional(),
	}
}

func (Problem) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("chat_id").
			Annotations(entsql.IndexWhere("(resolved_at IS NULL AND manager_id IS NULL)")).Unique(),
		index.Fields("chat_id").
			Annotations(entsql.IndexWhere("(resolved_at IS NULL AND manager_id IS NOT NULL)")).Unique().
			StorageKey("problems_chat_id_idx"),
	}
}

// Edges of the Problem.
func (Problem) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("chat", Chat.Type).Required().Field("chat_id").Immutable().Unique(),
		edge.From("messages", Message.Type).Ref("problem"),
	}
}
