package main

import (
	"fmt"

	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	"github.com/keepcalmist/chat-service/internal/store"
	gethistory "github.com/keepcalmist/chat-service/internal/usecases/client/get-history"
	"go.uber.org/zap"

	keycloakclient "github.com/keepcalmist/chat-service/internal/clients/keycloak"
	"github.com/keepcalmist/chat-service/internal/config"
	server_client "github.com/keepcalmist/chat-service/internal/server-client"
	clientv1 "github.com/keepcalmist/chat-service/internal/server-client/v1"
)

const nameServerClient = "server-client"

func initServerClient(
	addr string,
	allowOrigins []string,
	role string,
	resource string,
	keycloakConfig config.Keycloak,
	database *store.Database,
	isProduction bool,
) (*server_client.Server, error) {
	lg := zap.L().Named(nameServerClient)

	repoMsg, err := messagesrepo.New(messagesrepo.NewOptions(
		database,
	))
	if err != nil {
		return nil, fmt.Errorf("init messages repo: %v", err)
	}

	getHistoryUsecase, err := gethistory.New(
		gethistory.NewOptions(repoMsg),
	)
	if err != nil {
		return nil, fmt.Errorf("init get history usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(
		clientv1.NewOptions(getHistoryUsecase, lg),
	)
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	swag, err := clientv1.GetSwagger()
	if err != nil {
		return nil, fmt.Errorf("get swagger: %v", err)
	}

	keyCloakClient, err := keycloakclient.New(
		keycloakclient.NewOptions(
			keycloakConfig.BasePath,
			keycloakConfig.Realm,
			keycloakConfig.ClientID,
			keycloakConfig.ClientSecret,
			keycloakclient.WithDebugMode(keycloakConfig.DebugMode),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("init keycloak client: %v", err)
	}

	srv, err := server_client.New(server_client.NewOptions(
		zap.L().Named(nameServerClient),
		addr,
		allowOrigins,
		swag,
		v1Handlers,
		keyCloakClient,
		role,
		resource,
		isProduction,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
