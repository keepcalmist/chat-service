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
	EventID   types.EventID   `validate:"required" json:"eventId"`
	RequestID types.RequestID `validate:"required" json:"requestId"`
	MessageID types.MessageID `validate:"required" json:"messageId"`
}

func (e MessageSentEvent) Validate() error { return validator.Validator.Struct(e) }

func NewMessageSentEvent(
	id types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) *MessageSentEvent {
	return &MessageSentEvent{
		EventID:   id,
		RequestID: requestID,
		MessageID: messageID,
	}
}

type MessageBlockedEvent struct {
	event

	EventID   types.EventID   `validate:"required"`
	RequestID types.RequestID `validate:"required"`
	MessageID types.MessageID `validate:"required"`
}

func (e MessageBlockedEvent) Validate() error { return validator.Validator.Struct(e) }

func NewMessageBlockedEvent(
	id types.EventID,
	requestID types.RequestID,
	messageID types.MessageID,
) *MessageSentEvent {
	return &MessageSentEvent{
		EventID:   id,
		RequestID: requestID,
		MessageID: messageID,
	}
}

type NewMessageEvent struct {
	event
	EventID   types.EventID   `json:"eventId,omitempty"  validate:"required" `
	RequestID types.RequestID `json:"requestId,omitempty"  validate:"required" `
	ChatID    types.ChatID    `json:"chatId,omitempty"  validate:"required" `
	MessageID types.MessageID `json:"messageId,omitempty"  validate:"required" `
	AuthorID  types.UserID    `json:"authorId,omitempty"  validate:"required_without=IsService" `
	CreatedAt time.Time       `json:"createdAt"  validate:"required" `
	Body      string          `json:"body,omitempty"  validate:"required" `
	IsService bool            `json:"isService"`
}

func (e NewMessageEvent) Validate() error { return validator.Validator.Struct(e) }

func NewNewMessageEvent(
	id types.EventID,
	requestID types.RequestID,
	chatID types.ChatID,
	messageID types.MessageID,
	authorID types.UserID,
	createdAt time.Time,
	messageBody string,
	IsService bool,
) *NewMessageEvent {
	return &NewMessageEvent{
		EventID:   id,
		RequestID: requestID,
		ChatID:    chatID,
		MessageID: messageID,
		AuthorID:  authorID,
		CreatedAt: createdAt,
		Body:      messageBody,
		IsService: IsService,
	}
}
