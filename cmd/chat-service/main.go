package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/keepcalmist/chat-service/internal/config"
	"github.com/keepcalmist/chat-service/internal/logger"
	serverdebug "github.com/keepcalmist/chat-service/internal/server-debug"
	"github.com/keepcalmist/chat-service/internal/store"
)

var configPath = flag.String("config", "configs/config.toml", "Path to config file")

func main() {
	if err := run(); err != nil {
		log.Fatalf("run app: %v", err)
	}
}

func run() (errReturned error) {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.ParseAndValidate(*configPath)
	if err != nil {
		return fmt.Errorf("parse and validate config %q: %v", *configPath, err)
	}

	setLevel, err := logger.Init(
		logger.NewOptions(cfg.Log.Level,
			cfg.Global.Env,
			logger.WithProductionMode(cfg.Global.IsProduction),
			logger.WithSentryDSN(cfg.Sentry.DSN),
		),
	)
	if err != nil {
		return fmt.Errorf("init logger error: %w", err)
	}
	defer logger.Sync()

	psqlClient, err := store.NewPSQLClient(store.NewPSQLOptions(
		cfg.Postgres.Address,
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Database,
		store.WithDebug(cfg.Postgres.Debug),
	))
	if err != nil {
		return fmt.Errorf("init db driver: %v", err)
	}
	defer func() {
		if err := psqlClient.Close(); err != nil {
			errReturned = fmt.Errorf("close db connection: %v", err)
		}
	}()

	if err = psqlClient.Schema.Create(ctx); err != nil {
		return fmt.Errorf("migrate db: %v", err)
	}

	srvDebug, err := serverdebug.New(
		serverdebug.NewOptions(
			cfg.Servers.Debug.Addr,
			serverdebug.WithLvlSetter(setLevel)),
	)
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	if cfg.Clients.Keycloak.DebugMode && cfg.Global.IsProduction() {
		zap.L().Warn("keycloak debug mode enabled in production")
	}

	srvClient, err := initServerClient(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		cfg.Servers.Client.RequiredAccess.Role,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Clients.Keycloak,
		store.NewDatabase(psqlClient),
		cfg.Global.IsProduction(),
	)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return srvClient.Run(ctx) })
	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
