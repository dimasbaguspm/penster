package observability

import (
	"context"
	"log/slog"

	"github.com/dimasbaguspm/penster/config"
)

type Observability struct {
	TracerShutdown func(context.Context) error
	MeterShutdown  func(context.Context) error
	Logger         *slog.Logger
}

func Init(ctx context.Context, cfg *config.Config) *Observability {
	obs := &Observability{}

	obs.Logger = InitLogger(ctx, cfg)
	obs.TracerShutdown = InitTracing(ctx, cfg)
	obs.MeterShutdown = InitMeter(ctx, cfg)

	slog.Info("observability initialized",
		"service", cfg.App.Env,
		"version", cfg.App.Version,
		"env", cfg.App.Env,
	)

	return obs
}

func (o *Observability) Shutdown(ctx context.Context) {
	if o.TracerShutdown != nil {
		o.TracerShutdown(ctx)
	}
	if o.MeterShutdown != nil {
		o.MeterShutdown(ctx)
	}
	o.Logger.Info("observability shutdown complete")
}
