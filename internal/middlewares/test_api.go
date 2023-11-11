package middlewares

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/keepcalmist/chat-service/internal/types"
	"github.com/labstack/echo/v4"
)

func SetToken(c echo.Context, uid types.UserID) {
	c.Set(tokenCtxKey, &jwt.Token{
		Claims: claims{
			Subject: uid,
		},
	})
}
