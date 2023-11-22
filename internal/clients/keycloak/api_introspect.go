package keycloakclient

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v4"
)

type IntrospectTokenResult struct {
	Exp    int              `json:"exp"`
	Iat    int              `json:"iat"`
	Aud    jwt.ClaimStrings `json:"aud"`
	Active bool             `json:"active"`
}

// IntrospectToken implements
// https://www.keycloak.org/docs/latest/authorization_services/index.html#obtaining-information-about-an-rpt
func (c *Client) IntrospectToken(ctx context.Context, token string) (*IntrospectTokenResult, error) {
	// http://${host}:${port}/realms/${realm_name}/protocol/openid-connect/token/introspect
	url := fmt.Sprintf("realms/%s/protocol/openid-connect/token/introspect", c.realm)
	var result IntrospectTokenResult
	resp, err := c.cli.R().SetContext(ctx).
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBasicAuth(c.clientID, c.secret).
		SetFormData(map[string]string{
			"token_type_hint": "requesting_party_token",
			"token":           token,
		}).SetResult(&result).
		Post(url)
	if err != nil {
		return nil, fmt.Errorf("send request to keycloak: %v", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("errored keycloak response: %v", resp.Status())
	}

	return &result, err
}

func (c *Client) auth(ctx context.Context) *resty.Request {
	return c.cli.R().SetContext(ctx)
}
