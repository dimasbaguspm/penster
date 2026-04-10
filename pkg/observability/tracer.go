package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dimasbaguspm/penster/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func InitTracer(ctx context.Context, cfg *config.Config) func(context.Context) {
	if !cfg.OTEL.Enabled {
		slog.Info("OTEL tracing is disabled")
		return func(context.Context) {}
	}

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		os.Exit(1)
		return nil
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("penster"),
			semconv.ServiceVersion(cfg.App.Version),
			attribute.String("deployment.environment", cfg.App.Env),
		),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
		os.Exit(1)
		return nil
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	otel.SetTracerProvider(tp)
	tracer = tp.Tracer("penster")

	slog.Info("OTEL tracer initialized", "service", "penster", "environment", cfg.App.Env)

	shutdown := func(ctx context.Context) {
		if err := tp.Shutdown(ctx); err != nil {
			slog.Error("failed to shutdown tracer provider", "error", err)
		}
	}

	return shutdown
}

func Tracer() trace.Tracer {
	return tracer
}
