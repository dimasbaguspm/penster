package observability

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/dimasbaguspm/penster/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var meter metric.Meter

// InitMeter initializes the OTEL meter with OTLP exporter to Mimir.
func InitMeter(ctx context.Context, cfg *config.Config) func(context.Context) error {
	if !cfg.Observability.Enabled {
		slog.Info("metrics disabled")
		return func(context.Context) error { return nil }
	}

	conn, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.Observability.MetricsEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create metrics exporter: %v", err))
		return func(context.Context) error { return err }
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.App.Env),
			semconv.ServiceVersion(cfg.App.Version),
			attribute.String("deployment.environment", cfg.App.Env),
		),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create resource: %v", err))
		return func(context.Context) error { return err }
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(conn)),
	)

	otel.SetMeterProvider(provider)
	meter = provider.Meter("penster")

	slog.Info("metrics initialized", "endpoint", cfg.Observability.MetricsEndpoint)

	return func(ctx context.Context) error {
		if err := provider.Shutdown(ctx); err != nil {
			slog.Error(fmt.Sprintf("failed to shutdown meter provider: %v", err))
			return err
		}
		return nil
	}
}

// Meter returns the global meter instance.
func Meter() metric.Meter {
	return meter
}