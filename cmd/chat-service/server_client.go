package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/keepcalmist/chat-service/internal/clients/keycloak"
	"github.com/keepcalmist/chat-service/internal/config"
	chatsrepo "github.com/keepcalmist/chat-service/internal/repositories/chats"
	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	problemsrepo "github.com/keepcalmist/chat-service/internal/repositories/problems"
	server_client "github.com/keepcalmist/chat-service/internal/server-client"
	clientv1 "github.com/keepcalmist/chat-service/internal/server-client/v1"
	"github.com/keepcalmist/chat-service/internal/services/outbox"
	"github.com/keepcalmist/chat-service/internal/store"
	gethistory "github.com/keepcalmist/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/keepcalmist/chat-service/internal/usecases/client/send-message"
)

const nameServerClient = "server-client"

func initServerClient(
	addr string,
	allowOrigins []string,
	role string,
	resource string,
	keycloakConfig config.Keycloak,
	isProduction bool,
	swag *openapi3.T,
	database *store.Database,
	chatRepository *chatsrepo.Repo,
	msgRepository *messagesrepo.Repo,
	problemRepository *problemsrepo.Repo,
	outboxService *outbox.Service,
) (*server_client.Server, error) {
	getHistoryUsecase, err := gethistory.New(
		gethistory.NewOptions(msgRepository),
	)
	if err != nil {
		return nil, fmt.Errorf("init get history usecase: %v", err)
	}

	sendMessageUsecase, err := sendmessage.New(
		sendmessage.NewOptions(
			chatRepository,
			msgRepository,
			outboxService,
			problemRepository,
			database,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("init send message usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(
		clientv1.NewOptions(getHistoryUsecase, sendMessageUsecase),
	)
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
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
