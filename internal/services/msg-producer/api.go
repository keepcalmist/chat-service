package msgproducer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"

	"github.com/keepcalmist/chat-service/internal/types"
)

type Message struct {
	ID         types.MessageID
	ChatID     types.ChatID
	Body       string
	FromClient bool
}

type transportMsg struct {
	ID         types.MessageID `json:"id"`
	ChatID     types.ChatID    `json:"chatId"`
	Body       string          `json:"body"`
	FromClient bool            `json:"fromClient"`
}

func (m Message) String() string {
	return fmt.Sprintf("Message{ID: %s, ChatID: %s, Body: %s, FromClient: %t}", m.ID, m.ChatID, m.Body, m.FromClient)
}

func (s *Service) ProduceMessage(ctx context.Context, msg Message) error {
	jsonData, err := json.Marshal(transportMsg(msg))
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	if s.cipher != nil {
		nonce, err := s.nonceFactory(s.cipher.NonceSize())
		if err != nil {
			return fmt.Errorf("failed to generate nonce: %w", err)
		}
		jsonData = s.cipher.Seal(nonce, nonce, jsonData, nil)
	}

	err = s.wr.WriteMessages(ctx, kafka.Message{
		Key:   []byte(msg.ChatID.String()),
		Value: jsonData,
	})
	if err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (s *Service) Close() error {
	return s.wr.Close()
}
