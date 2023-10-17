package main

import (
	"fmt"

	"github.com/keepcalmist/chat-service/internal/server-client"
	"go.uber.org/zap"

	clientv1 "github.com/keepcalmist/chat-service/internal/server-client/v1"
)

const nameServerClient = "server-client"

func initServerClient( // FIXME: воспользуйся мной в chat-service/main.go
	addr string,
	allowOrigins []string,
	swaggerFile string,
) (*server_client.Server, error) {
	lg := zap.L().Named(nameServerClient)

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(lg))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	swag, err := clientv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get swagger: %v", err)
	}

	srv, err := server_client.New(server_client.NewOptions(
		zap.L().Named(nameServerClient),
		addr,
		allowOrigins,
		swag,
		v1Handlers,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
