package middlewares

import (
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v4"

	"github.com/keepcalmist/chat-service/internal/types"
)

var (
	ErrNoAllowedResources = errors.New("no allowed resources")
	ErrSubjectNotDefined  = errors.New(`"sub" is not defined`)
)

type claims struct {
	jwt.RegisteredClaims
	Subject        types.UserID                      `json:"sub"`
	ResourceAccess map[string]map[string]interface{} `json:"resource_access"` //nolint
}

// Valid returns errors:
// - from StandardClaims validation;
// - ErrNoAllowedResources, if claims doesn't contain `resource_access` map or it's empty;
// - ErrSubjectNotDefined, if claims doesn't contain `sub` field or subject is zero UUID.
func (c claims) Valid() error {
	if err := c.RegisteredClaims.Valid(); err != nil {
		return err
	}

	if c.Subject == types.UserIDNil {
		return ErrSubjectNotDefined
	}
	if len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}

	if len(c.ResourceAccess) == 0 {
		return ErrNoAllowedResources
	}

	return nil
}

func (c claims) UserID() types.UserID {
	return c.Subject
}

func parse(tokenString string) (*jwt.Token, error) {
	token, _, err := jwt.NewParser().ParseUnverified(tokenString, &claims{})
	if err != nil {
		return nil, fmt.Errorf("failed to parse claims: %w", err)
	}
	return token, err
}
