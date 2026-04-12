package observability

import (
	"context"
	"log/slog"
	"os"

	"github.com/dimasbaguspm/penster/config"
)

var logger *slog.Logger

func InitLogger(ctx context.Context, cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)

	if cfg.Observability.Enabled {
		slog.Info("logger initialized",
			"service", cfg.App.Env,
			"env", cfg.App.Env,
		)
	}

	return logger
}

func Logger() *slog.Logger {
	return logger
}
