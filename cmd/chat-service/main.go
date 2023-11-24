package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	clientevents "github.com/keepcalmist/chat-service/internal/server/server-client/events"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/keepcalmist/chat-service/internal/config"
	"github.com/keepcalmist/chat-service/internal/logger"
	chatsrepo "github.com/keepcalmist/chat-service/internal/repositories/chats"
	jobsrepo "github.com/keepcalmist/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/keepcalmist/chat-service/internal/repositories/messages"
	problemsrepo "github.com/keepcalmist/chat-service/internal/repositories/problems"
	serverdebug "github.com/keepcalmist/chat-service/internal/server-debug"
	clientv1 "github.com/keepcalmist/chat-service/internal/server/server-client/v1"
	managerv1 "github.com/keepcalmist/chat-service/internal/server/server-manager/v1"
	managerload "github.com/keepcalmist/chat-service/internal/services/manager-load"
	inmemmanagerpool "github.com/keepcalmist/chat-service/internal/services/manager-pool/in-mem"
	msgproducer "github.com/keepcalmist/chat-service/internal/services/msg-producer"
	"github.com/keepcalmist/chat-service/internal/store"
	"github.com/keepcalmist/chat-service/pkg/shutdown"
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

	clientSwagger, err := clientv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get client swagger: %v", err)
	}

	managerSwagger, err := managerv1.GetSwagger()
	if err != nil {
		return fmt.Errorf("get manager swagger: %v", err)
	}
	eventsSchema, err := clientevents.GetSwagger()
	if err != nil {
		return fmt.Errorf("get events swagger: %v", err)
	}

	srvDebug, err := serverdebug.New(
		serverdebug.NewOptions(
			cfg.Servers.Debug.Addr,
			clientSwagger,
			managerSwagger,
			eventsSchema,
			serverdebug.WithLvlSetter(setLevel)),
	)
	if err != nil {
		return fmt.Errorf("init debug server: %v", err)
	}

	if cfg.Clients.Keycloak.DebugMode && cfg.Global.IsProduction() {
		zap.L().Warn("keycloak debug mode enabled in production")
	}

	if cfg.Services.MsgProducer.EncryptKey == "" {
		zap.L().Warn("msg producer encrypt disabled")
	}

	database := store.NewDatabase(psqlClient)

	repoMsg, err := messagesrepo.New(messagesrepo.NewOptions(
		database,
	))
	if err != nil {
		return fmt.Errorf("init messages repo: %v", err)
	}

	repoChat, err := chatsrepo.New(chatsrepo.NewOptions(
		database,
	))
	if err != nil {
		return fmt.Errorf("init chats repo: %v", err)
	}

	repoProblems, err := problemsrepo.New(problemsrepo.NewOptions(
		database,
	))
	if err != nil {
		return fmt.Errorf("init problems repo: %v", err)
	}

	repoJobs, err := jobsrepo.New(jobsrepo.NewOptions(
		database,
	))
	if err != nil {
		return fmt.Errorf("init jobs repo: %v", err)
	}

	kafkaWriter := msgproducer.NewKafkaWriter(
		cfg.Services.MsgProducer.Brokers,
		cfg.Services.MsgProducer.Topic,
		cfg.Services.MsgProducer.BatchSize,
	)

	producer, err := msgproducer.New(
		msgproducer.NewOptions(
			kafkaWriter,
			msgproducer.WithEncryptKey(cfg.Services.MsgProducer.EncryptKey),
		),
	)
	if err != nil {
		return fmt.Errorf("init msg producer: %v", err)
	}

	managerLoadService, err := managerload.New(managerload.NewOptions(cfg.Services.ManagerLoad.MaxProblemsAtSameTime, repoProblems))
	if err != nil {
		return fmt.Errorf("init manager load service: %v", err)
	}

	poolService := inmemmanagerpool.New()

	shutdownChan := shutdown.NewShutDown()
	outboxService, err := initOutbox(cfg.Services, database, repoJobs, repoMsg, producer)
	if err != nil {
		return fmt.Errorf("init outbox: %v", err)
	}

	srvManager, err := initServerManager(
		cfg.Servers.Manager.Addr,
		cfg.Servers.Manager.AllowOrigins,
		cfg.Servers.Manager.RequiredAccess.Role,
		cfg.Servers.Manager.RequiredAccess.Resource,
		cfg.Clients.Keycloak,
		cfg.Global.IsProduction(),
		managerSwagger,
		managerLoadService,
		poolService,
		shutdownChan,
		cfg.Servers.Manager.SecWSProtocol,
	)
	if err != nil {
		return fmt.Errorf("init manager server: %v", err)
	}

	srvClient, err := initServerClient(
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		cfg.Servers.Client.RequiredAccess.Role,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Clients.Keycloak,
		cfg.Global.IsProduction(),
		clientSwagger,
		database,
		repoChat,
		repoMsg,
		repoProblems,
		outboxService,
		shutdownChan,
		cfg.Servers.Manager.SecWSProtocol,
	)
	if err != nil {
		return fmt.Errorf("init client server: %v", err)
	}

	eg, ctx := errgroup.WithContext(ctx)
	// Run servers.
	eg.Go(func() error { return srvDebug.Run(ctx) })
	eg.Go(func() error { return srvClient.Run(ctx) })
	eg.Go(func() error { return srvManager.Run(ctx) })
	eg.Go(func() error { return outboxService.Run(ctx) })
	// Run services.
	// Ждут своего часа.
	// ...

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}

	return nil
}
