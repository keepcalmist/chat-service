package middlewares

import (
	"github.com/keepcalmist/chat-service/internal/types"
	"github.com/labstack/echo/v4"
)

func SetToken(c echo.Context, uid types.UserID) {
	// FIXME: В контекст по ключу tokenCtxKey необходимо положить jwt.Token с клэймсами:
	// FIXME: - которые всегда валидные
	// FIXME: - из которых можно достать uid
}
