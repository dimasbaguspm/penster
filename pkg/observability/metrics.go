package observability

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/dimasbaguspm/penster/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

var meter metric.Meter

// Metric instruments organized by category
var (
	// Business metrics - transaction operations
	TransactionsCreated metric.Int64Counter
	TransactionsUpdated metric.Int64Counter
	TransactionsDeleted metric.Int64Counter
	DraftsConfirmed     metric.Int64Counter
	DraftsRejected      metric.Int64Counter

	// Business metrics - account operations
	AccountsCreated metric.Int64Counter
	AccountsUpdated metric.Int64Counter
	AccountsDeleted metric.Int64Counter

	// Business metrics - category operations
	CategoriesCreated metric.Int64Counter
	CategoriesUpdated metric.Int64Counter
	CategoriesDeleted metric.Int64Counter

	// Infrastructure metrics - scheduler
	SchedulerJobsExecuted metric.Int64Counter
	SchedulerJobsFailed   metric.Int64Counter
	SchedulerJobDuration  metric.Float64Histogram

	// Traffic metrics - HTTP
	HTTPRequestsTotal   metric.Int64Counter
	HTTPRequestDuration metric.Float64Histogram
	HTTPPanicCount      metric.Int64Counter
)

// InitMeter initializes the OTEL meter with OTLP exporter to Mimir.
func InitMeter(ctx context.Context, cfg *config.Config) func(context.Context) error {
	// Always register no-op metrics first so Add() calls are safe
	registerMetrics(noop.NewMeterProvider().Meter("penster"))

	if !cfg.Observability.Enabled {
		slog.Info("metrics disabled, using no-op instruments")
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

	registerMetrics(meter)

	slog.Info("metrics initialized", "endpoint", cfg.Observability.MetricsEndpoint)

	return func(ctx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(ctx, 2*time.Second)
		defer cancel()
		if err := provider.Shutdown(shutdownCtx); err != nil {
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

// RecordJobMetrics records scheduler job execution metrics.
func RecordJobMetrics(ctx context.Context, jobName string, success bool, durationMs float64) {
	if success {
		SchedulerJobsExecuted.Add(ctx, 1)
	} else {
		SchedulerJobsFailed.Add(ctx, 1)
	}
	SchedulerJobDuration.Record(ctx, durationMs)
}

// registerMetrics creates and registers all metric instruments.
func registerMetrics(m metric.Meter) error {
	// Business metrics - transactions
	var err error
	TransactionsCreated, err = m.Int64Counter("penster_transactions_created",
		metric.WithDescription("Total number of transactions created"),
		metric.WithUnit("{transaction}"))
	if err != nil {
		return fmt.Errorf("failed to create transactions_created counter: %w", err)
	}
	TransactionsUpdated, err = m.Int64Counter("penster_transactions_updated",
		metric.WithDescription("Total number of transactions updated"),
		metric.WithUnit("{transaction}"))
	if err != nil {
		return fmt.Errorf("failed to create transactions_updated counter: %w", err)
	}
	TransactionsDeleted, err = m.Int64Counter("penster_transactions_deleted",
		metric.WithDescription("Total number of transactions deleted"),
		metric.WithUnit("{transaction}"))
	if err != nil {
		return fmt.Errorf("failed to create transactions_deleted counter: %w", err)
	}
	DraftsConfirmed, err = m.Int64Counter("penster_drafts_confirmed",
		metric.WithDescription("Total number of drafts confirmed"),
		metric.WithUnit("{draft}"))
	if err != nil {
		return fmt.Errorf("failed to create drafts_confirmed counter: %w", err)
	}
	DraftsRejected, err = m.Int64Counter("penster_drafts_rejected",
		metric.WithDescription("Total number of drafts rejected"),
		metric.WithUnit("{draft}"))
	if err != nil {
		return fmt.Errorf("failed to create drafts_rejected counter: %w", err)
	}

	// Business metrics - accounts
	AccountsCreated, err = m.Int64Counter("penster_accounts_created",
		metric.WithDescription("Total number of accounts created"),
		metric.WithUnit("{account}"))
	if err != nil {
		return fmt.Errorf("failed to create accounts_created counter: %w", err)
	}
	AccountsUpdated, err = m.Int64Counter("penster_accounts_updated",
		metric.WithDescription("Total number of accounts updated"),
		metric.WithUnit("{account}"))
	if err != nil {
		return fmt.Errorf("failed to create accounts_updated counter: %w", err)
	}
	AccountsDeleted, err = m.Int64Counter("penster_accounts_deleted",
		metric.WithDescription("Total number of accounts deleted"),
		metric.WithUnit("{account}"))
	if err != nil {
		return fmt.Errorf("failed to create accounts_deleted counter: %w", err)
	}

	// Business metrics - categories
	CategoriesCreated, err = m.Int64Counter("penster_categories_created",
		metric.WithDescription("Total number of categories created"),
		metric.WithUnit("{category}"))
	if err != nil {
		return fmt.Errorf("failed to create categories_created counter: %w", err)
	}
	CategoriesUpdated, err = m.Int64Counter("penster_categories_updated",
		metric.WithDescription("Total number of categories updated"),
		metric.WithUnit("{category}"))
	if err != nil {
		return fmt.Errorf("failed to create categories_updated counter: %w", err)
	}
	CategoriesDeleted, err = m.Int64Counter("penster_categories_deleted",
		metric.WithDescription("Total number of categories deleted"),
		metric.WithUnit("{category}"))
	if err != nil {
		return fmt.Errorf("failed to create categories_deleted counter: %w", err)
	}

	// Infrastructure metrics - scheduler
	SchedulerJobsExecuted, err = m.Int64Counter("penster_scheduler_jobs_executed",
		metric.WithDescription("Total number of scheduler jobs executed"),
		metric.WithUnit("{job}"))
	if err != nil {
		return fmt.Errorf("failed to create scheduler_jobs_executed counter: %w", err)
	}
	SchedulerJobsFailed, err = m.Int64Counter("penster_scheduler_jobs_failed",
		metric.WithDescription("Total number of scheduler jobs failed"),
		metric.WithUnit("{job}"))
	if err != nil {
		return fmt.Errorf("failed to create scheduler_jobs_failed counter: %w", err)
	}
	SchedulerJobDuration, err = m.Float64Histogram("penster_scheduler_job_duration",
		metric.WithDescription("Duration of scheduler jobs"),
		metric.WithUnit("ms"))
	if err != nil {
		return fmt.Errorf("failed to create scheduler_job_duration histogram: %w", err)
	}

	// Traffic metrics - HTTP
	HTTPRequestsTotal, err = m.Int64Counter("penster_http_requests_total",
		metric.WithDescription("Total number of HTTP requests"),
		metric.WithUnit("{request}"))
	if err != nil {
		return fmt.Errorf("failed to create http_requests_total counter: %w", err)
	}
	HTTPRequestDuration, err = m.Float64Histogram("penster_http_request_duration",
		metric.WithDescription("Duration of HTTP requests"),
		metric.WithUnit("ms"))
	if err != nil {
		return fmt.Errorf("failed to create http_request_duration histogram: %w", err)
	}
	HTTPPanicCount, err = m.Int64Counter("penster_http_panic_count",
		metric.WithDescription("Total number of HTTP panics recovered"),
		metric.WithUnit("{panic}"))
	if err != nil {
		return fmt.Errorf("failed to create http_panic_count counter: %w", err)
	}

	return nil
}
