package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"github.com/TheZeroSlave/zapsentry"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/keepcalmist/chat-service/internal/buildinfo"
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode func() bool
	sentryDSN      string
	env            string `option:"mandatory" validate:"required,oneof=dev stage prod"`
}

func MustInit(opts Options) func(level zapcore.Level) {
	f, err := Init(opts)
	if err != nil {
		panic(err)
	}
	return f
}

func Init(opts Options) (func(level zapcore.Level), error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	lvl, err := zap.ParseAtomicLevel(opts.level)
	if err != nil {
		return nil, err
	}
	var enc zapcore.Encoder
	cfg := zapcore.EncoderConfig{
		TimeKey:    "T",
		LevelKey:   "level",
		MessageKey: "msg",
		NameKey:    "component",
		EncodeTime: zapcore.ISO8601TimeEncoder,
	}

	if opts.productionMode != nil && opts.productionMode() {
		cfg.EncodeLevel = zapcore.CapitalLevelEncoder
		enc = zapcore.NewJSONEncoder(cfg)
	} else {
		cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
		enc = zapcore.NewConsoleEncoder(cfg)
	}

	cores := []zapcore.Core{
		zapcore.NewCore(enc, os.Stdout, lvl),
	}
	if opts.sentryDSN != "" {
		cl, err := NewSentryClient(opts.sentryDSN, opts.env, buildinfo.BuildInfo.Main.Version)
		if err != nil {
			return nil, fmt.Errorf("new sentry client: %v", err)
		}
		sentryCore, err := zapsentry.NewCore(zapsentry.Configuration{
			Level: zap.WarnLevel,
			FrameMatcher: zapsentry.CombineFrameMatchers(
				zapsentry.SkipFunctionPrefixFrameMatcher("go.uber.org/zap"),
			),
		}, zapsentry.NewSentryClientFromClient(cl))
		if err != nil {
			return nil, fmt.Errorf("new sentry core: %v", err)
		}
		cores = append(cores, sentryCore)
	}
	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return lvl.SetLevel, nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
