package canreceiveproblems

import (
	"github.com/keepcalmist/chat-service/internal/types"
	"github.com/keepcalmist/chat-service/internal/validator"
)

type Request struct {
	ID        types.RequestID `validate:"required"`
	ManagerID types.UserID    `validate:"required"`
}

func (r Request) Validate() error {
	if err := validator.Validator.Struct(r); err != nil {
		return err
	}

	return nil
}

type Response struct {
	Result bool
}
