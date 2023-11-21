package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/keepcalmist/chat-service/internal/clients/keycloak"
	"github.com/keepcalmist/chat-service/internal/config"
	"github.com/keepcalmist/chat-service/internal/server"
	managerv1 "github.com/keepcalmist/chat-service/internal/server/server-manager/v1"
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
) (*server.Server, error) {
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

	srv, err := server.New(server.NewOptions(
		zap.L().Named(nameServerManager),
		addr,
		allowOrigins,
		swag,
		keyCloakClient,
		role,
		resource,
		isProduction,
	), managerv1.Handlers{}) // TODO реализовать, как только будет возможно
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
