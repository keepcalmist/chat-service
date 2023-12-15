package afcverdictsprocessor

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"

	"github.com/dgrijalva/jwt-go"

	"github.com/keepcalmist/chat-service/internal/types"
	"github.com/keepcalmist/chat-service/internal/validator"
)

const (
	OK         VerdictStatus = "ok"
	Suspicious VerdictStatus = "suspicious"
)

type VerdictStatus string

type verdict struct {
	ChatID    types.ChatID    `json:"chatId" validate:"required"`
	MessageID types.MessageID `json:"messageId" validate:"required"`
	Status    VerdictStatus   `json:"status" validate:"oneof=ok suspicious"`
}

func unmarshalVerdictWithKey(pubKey *rsa.PublicKey, data []byte) (*verdict, error) {
	value, err := jwt.Parse(string(data), func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", jwtToken.Header["alg"])
		}

		return pubKey, nil
	})
	if err != nil {
		return nil, fmt.Errorf("validate jwt-token: %w", err)
	}

	claims, ok := value.Claims.(jwt.MapClaims)
	if !ok || !value.Valid {
		return nil, fmt.Errorf("validate token: invalid")
	}

	unmarshalledVerdict, err := unmarshalVerdictFromClaims(claims)
	if err != nil {
		return nil, fmt.Errorf("unmarshal verdict: %w", err)
	}

	if err := validator.Validator.Struct(unmarshalledVerdict); err != nil {
		return nil, fmt.Errorf("validate verdict: %w", err)
	}

	return unmarshalledVerdict, nil
}

func unmarshalVerdictFromClaims(data jwt.MapClaims) (*verdict, error) {
	value, ok := data["chatId"]
	if !ok {
		return nil, fmt.Errorf("chatId is invalid")
	}

	chatID := &types.ChatID{}
	err := chatID.Scan(value)
	if err != nil || chatID.IsZero() {
		return nil, fmt.Errorf("chatId is invalid")
	}

	value, ok = data["messageId"]
	if !ok {
		return nil, fmt.Errorf("messageId is invalid")
	}
	messageID := &types.MessageID{}
	err = messageID.Scan(value)
	if err != nil || messageID.IsZero() {
		return nil, fmt.Errorf("messageId is invalid")
	}

	status, ok := data["status"]
	if !ok {
		return nil, fmt.Errorf("status is invalid")
	}

	transformedStatus, ok := status.(string)
	if !ok {
		return nil, fmt.Errorf("status is invalid")
	}

	return &verdict{
		ChatID:    *chatID,
		MessageID: *messageID,
		Status:    VerdictStatus(transformedStatus),
	}, nil
}

func unmarshalVerdictFromJSON(data []byte) (*verdict, error) {
	var v verdict
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("unmarshal verdict: %w", err)
	}

	if err := validator.Validator.Struct(v); err != nil {
		return nil, fmt.Errorf("validate verdict: %w", err)
	}

	return &v, nil
}
