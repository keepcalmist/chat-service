package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/keepcalmist/chat-service/internal/middlewares"
	"github.com/keepcalmist/chat-service/internal/server/server-client/errhandler"
	clientv1 "github.com/keepcalmist/chat-service/internal/server/server-client/v1"
	managerv1 "github.com/keepcalmist/chat-service/internal/server/server-manager/v1"
	"github.com/keepcalmist/chat-service/pkg/shutdown"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

type handlersToRegister interface {
	clientv1.Handlers | managerv1.Handlers
}

type wsHandler interface {
	Serve(eCtx echo.Context) error
}

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger       *zap.Logger              `option:"mandatory" validate:"required"`
	addr         string                   `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins []string                 `option:"mandatory" validate:"min=1"`
	v1Swagger    *openapi3.T              `option:"mandatory" validate:"required"`
	introspector middlewares.Introspector `option:"mandatory" validate:"required"`
	role         string                   `option:"mandatory" validate:"required"`
	resource     string                   `option:"mandatory" validate:"required"`
	isProduction bool                     `option:"mandatory"`
	wsHandler    wsHandler                `option:"mandatory" validate:"required"`
	wsShutdown   *shutdown.ShutDown       `option:"mandatory" validate:"required"`
}

type Server struct {
	lg         *zap.Logger
	srv        *http.Server
	wsShutdown *shutdown.ShutDown
}

func New[t handlersToRegister](opts Options, handlers t) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	errHandle, err := errhandler.New(errhandler.NewOptions(
		opts.logger,
		opts.isProduction,
		errhandler.ResponseBuilder,
	))
	if err != nil {
		return nil, fmt.Errorf("create err handler: %w", err)
	}

	e := echo.New()

	e.HTTPErrorHandler = errHandle.Handle
	e.Use(
		middlewares.NewRecovery(opts.logger),
		middlewares.NewRequestLogger(opts.logger),
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: opts.allowOrigins,
			AllowMethods: []string{http.MethodPost},
		}),
		middleware.BodyLimit("3K"),
	)

	v1 := e.Group("v1",
		middlewares.NewKeycloakTokenAuth(opts.introspector, opts.resource, opts.role, false),
		oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
			Options: openapi3filter.Options{
				ExcludeRequestBody:  false,
				ExcludeResponseBody: true,
				AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
			},
		}))

	ws := e.Group("ws", middlewares.NewKeycloakTokenAuth(opts.introspector, opts.resource, opts.role, true))
	ws.GET("", opts.wsHandler.Serve)

	switch any(handlers).(type) {
	case clientv1.Handlers:
		clientv1.RegisterHandlers(v1, any(handlers).(clientv1.Handlers))
	case managerv1.Handlers:
		managerv1.RegisterHandlers(v1, any(handlers).(managerv1.Handlers))
	default:
		return nil, fmt.Errorf("handlers type is not supported")
	}

	return &Server{
		lg: opts.logger,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		wsShutdown: opts.wsShutdown,
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		err := s.wsShutdown.Close()
		if err != nil {
			s.lg.Error("failed to close ws shutdown", zap.Error(err))
		}

		return s.srv.Shutdown(ctx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}
