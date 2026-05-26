package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"foodtosave-notify/internal/app"
	"foodtosave-notify/internal/config"
)

func main() {
	os.Exit(run())
}

func run() int {
	configFlag := flag.String("config", "", "path to config.yaml")
	flag.Parse()

	cfgPath := strings.TrimSpace(*configFlag)
	if cfgPath == "" {
		cfgPath = strings.TrimSpace(os.Getenv("CONFIG_PATH"))
	}
	if cfgPath == "" {
		cfgPath = config.DefaultConfigPath
	}

	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo})
	logger := slog.New(handler)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		logger.Error("failed to load config", "path", cfgPath, "err", err)
		return 1
	}

	logger.Info("config loaded", "path", cfgPath, "interval", cfg.IntervalDuration().String())

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err = app.Run(sigCtx, cfg, logger)
	if err != nil && !errors.Is(err, context.Canceled) {
		logger.Error("execution terminated with error", "err", err)
		return 1
	}
	return 0
}
