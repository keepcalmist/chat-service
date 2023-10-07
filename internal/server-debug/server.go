package serverdebug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/pprof"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/sync/errgroup"

	"github.com/keepcalmist/chat-service/internal/buildinfo"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr      string `option:"mandatory" validate:"required,hostname_port"`
	logSetter func(level zapcore.Level)
}

type Server struct {
	lg        *zap.Logger
	srv       *http.Server
	lvlSetter func(level zapcore.Level)
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	lg := zap.L().Named("server-debug")

	e := echo.New()
	e.Use(middleware.Recover())

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
		lvlSetter: opts.logSetter,
	}
	index := newIndexPage()

	e.GET("/version", s.Version)
	index.addPage("/version", "Get build information")
	index.addPage("/debug/pprof/", "Go stg profile")
	index.addPage("/debug/pprof/profile?seconds=30", "Takes half-minute profile")

	// Обработка "/log/level"
	e.PUT("/log/level", s.SetLogLvl)

	// Обработка "/debug/pprof/" и связанных команд
	pprof.Register(e)

	e.GET("/", index.handler)
	return s, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

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

func (s *Server) Version(eCtx echo.Context) error {
	return eCtx.JSON(http.StatusOK, buildinfo.BuildInfo)
}

func (s *Server) SetLogLvl(eCtx echo.Context) error {
	req := eCtx.Request()

	lvl, err := zapcore.ParseLevel(req.FormValue("level"))
	if err != nil {
		return err
	}
	old := s.lg.Level().String()
	s.lvlSetter(lvl)

	s.lg.Info("switching log lvl",
		zap.String("old", old),
		zap.String("new", s.lg.Level().String()),
	)

	return nil
}
