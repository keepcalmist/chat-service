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
	v := new(verdict)
	chatID, ok := data["chatId"]
	if !ok {
		return nil, fmt.Errorf("chatId is invalid")
	}
	v.ChatID, ok = chatID.(types.ChatID)
	if !ok || v.ChatID == types.ChatIDNil {
		return nil, fmt.Errorf("chatId is invalid")
	}

	messageID, ok := data["messageId"]
	if !ok {
		return nil, fmt.Errorf("messageId is invalid")
	}

	v.MessageID, ok = messageID.(types.MessageID)
	if !ok || v.MessageID == types.MessageIDNil {
		return nil, fmt.Errorf("messageId is invalid")
	}

	status, ok := data["status"]
	if !ok {
		return nil, fmt.Errorf("status is invalid")
	}
	v.Status, ok = status.(VerdictStatus)
	if !ok {
		return nil, fmt.Errorf("status is invalid")
	}

	return v, nil
}

func unmarshalVerdictFromJSON(data []byte) (*verdict, error) {
	var v *verdict
	if err := json.Unmarshal(data, v); err != nil {
		return nil, fmt.Errorf("unmarshal verdict: %w", err)
	}

	if err := validator.Validator.Struct(v); err != nil {
		return nil, fmt.Errorf("validate verdict: %w", err)
	}

	return v, nil
}
