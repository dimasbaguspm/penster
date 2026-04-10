package observability

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/dimasbaguspm/penster/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracer trace.Tracer

func InitTracer(ctx context.Context, cfg *config.Config) func(context.Context) {
	if !cfg.OTEL.Enabled {
		slog.Info("OTEL tracing is disabled")
		tracer = noop.NewTracerProvider().Tracer("")
		return func(context.Context) {}
	}

	conn, err := grpc.NewClient(cfg.OTEL.Endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create gRPC connection to OTEL exporter: %v", err))
		os.Exit(1)
		return nil
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		slog.Error(fmt.Sprintf("failed to create OTLP trace exporter: %v", err))
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

	slog.Info("OTEL tracer initialized", "service", "penster", "environment", cfg.App.Env, "endpoint", cfg.OTEL.Endpoint)

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
