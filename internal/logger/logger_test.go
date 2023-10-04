package logger_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/keepcalmist/chat-service/internal/logger"
)

func TestInit(t *testing.T) {
	setter, err := logger.Init(logger.NewOptions("error", logger.WithProductionMode(func() bool {
		return true
	})))
	require.NoError(t, err)
	require.NotNil(t, setter)
	zap.L().Named("user-cache").Error("inconsistent state", zap.String("uid", "1234"))
	// {"level":"ERROR","T":"2022-10-09T13:56:47.626+0300","component":"user-cache","msg":"inconsistent state","uid":"1234"}
}
