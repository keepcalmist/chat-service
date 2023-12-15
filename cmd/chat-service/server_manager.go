package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/keepcalmist/chat-service/internal/clients/keycloak"
	"github.com/keepcalmist/chat-service/internal/config"
	"github.com/keepcalmist/chat-service/internal/server"
	managerv1 "github.com/keepcalmist/chat-service/internal/server/server-manager/v1"
	managerload "github.com/keepcalmist/chat-service/internal/services/manager-load"
	managerpool "github.com/keepcalmist/chat-service/internal/services/manager-pool"
	canreceiveproblems "github.com/keepcalmist/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/keepcalmist/chat-service/internal/usecases/manager/free-hands"
	websocketstream "github.com/keepcalmist/chat-service/internal/websocket-stream"
	"github.com/keepcalmist/chat-service/pkg/shutdown"
)

const nameServerManager = "server-manager"

func initServerManager(
	addr string,
	allowOrigins []string,
	role string,
	resource string,
	keycloakConfig config.Keycloak,
	isProduction bool,
	swag *openapi3.T,
	managerLoadService *managerload.Service,
	managerPoolService managerpool.Pool,
	shutdownChan *shutdown.ShutDown,
	secWsProtocol string,
) (*server.Server, error) {
	wsHandler, err := websocketstream.NewHTTPHandler(
		websocketstream.NewOptions(
			zap.L(),
			dummyEventStream{},
			dummyAdapter{},
			websocketstream.JSONEventWriter{},
			websocketstream.NewUpgrader(allowOrigins, secWsProtocol),
			shutdownChan.Done(),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("init ws handler: %v", err)
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

	useCaseCanReceiveProblem, err := canreceiveproblems.New(
		canreceiveproblems.NewOptions(managerLoadService, managerPoolService),
	)
	if err != nil {
		return nil, fmt.Errorf("init usecase can reciev problem: %v", err)
	}

	useCaseFreeHands, err := freehands.New(freehands.NewOptions(managerLoadService, managerPoolService))
	if err != nil {
		return nil, fmt.Errorf("init usecase free hands: %v", err)
	}

	handlers, err := managerv1.NewHandlers(
		managerv1.NewOptions(useCaseCanReceiveProblem, useCaseFreeHands),
	)
	if err != nil {
		return nil, fmt.Errorf("init handlers: %v", err)
	}

	srv, err := server.New(server.NewOptions(
		zap.L().Named(nameServerManager),
		addr,
		allowOrigins,
		swag,
		keyCloakClient,
		role,
		resource,
		isProduction,
		wsHandler,
		shutdownChan,
	), handlers)
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
