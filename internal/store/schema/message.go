package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/keepcalmist/chat-service/internal/types"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).Default(types.NewMessageID).Immutable(),
		field.UUID("author_id", types.UserID{}).Immutable().Optional(),
		field.Bool("is_visible_for_client").Default(true),
		field.Bool("is_visible_for_manager").Default(true),
		field.Text("body").NotEmpty().Immutable().MaxLen(1024),
		field.Time("checked_at").Nillable().Optional(),
		field.Bool("is_blocked").Optional(),
		field.Bool("is_service"),
		field.Time("created_at").Immutable().Default(time.Now),
		field.UUID("chat_id", types.ChatID{}).Immutable(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("chat", Chat.Type).Required().Unique().Field("chat_id").Immutable(),
		edge.To("problem", Problem.Type).Required().Unique(),
	}
}
