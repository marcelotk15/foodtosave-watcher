package app

import (
	"context"
	"log/slog"
	"time"

	"foodtosave-notify/internal/config"
	"foodtosave-notify/internal/watcher"
)

func Run(ctx context.Context, cfg *config.Config, log *slog.Logger) error {
	if log == nil {
		log = slog.Default()
	}

	w := watcher.New(cfg, log)
	w.Tick(ctx)

	ticker := time.NewTicker(cfg.IntervalDuration())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info("stopping watcher", "reason", ctx.Err())
			return ctx.Err()
		case <-ticker.C:
			w.Tick(ctx)
		}
	}
}
