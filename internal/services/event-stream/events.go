package eventstream

import (
	"time"

	"github.com/keepcalmist/chat-service/internal/types"
	"github.com/keepcalmist/chat-service/internal/validator"
)

type Event interface {
	eventMarker()
	Validate() error
}

type event struct{}         //
func (*event) eventMarker() {}

// MessageSentEvent indicates that the message was checked by AFC
// and was sent to the manager. Two gray ticks.
type MessageSentEvent struct {
	event
}

func (e MessageSentEvent) Validate() error { return nil }

type NewMessageEvent struct {
	event
	ID          types.EventID   `validate:"required"`
	RequestID   types.RequestID `validate:"required"`
	ChatID      types.ChatID    `validate:"required"`
	MessageID   types.MessageID `validate:"required"`
	UserID      types.UserID    `validate:"required"`
	CreatedAt   time.Time       `validate:"required"`
	MessageBody string          `validate:"required"`
	IsChecked   bool
}

func (e NewMessageEvent) Validate() error {
	return validator.Validator.Struct(e)
}

func NewNewMessageEvent(
	id types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	messageID types.MessageID,
	userID types.UserID,
	createdAt time.Time,
	messageBody string,
	IsChecked bool,
) *NewMessageEvent {
	return &NewMessageEvent{
		ID:          id,
		RequestID:   requestID,
		ChatID:      chatID,
		MessageID:   messageID,
		UserID:      userID,
		CreatedAt:   createdAt,
		MessageBody: messageBody,
		IsChecked:   IsChecked,
	}
}
