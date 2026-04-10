package observability

import (
	"context"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// StartServiceSpan starts a new span for a service method.
func StartServiceSpan(ctx context.Context, serviceName, operation string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, serviceName+"."+operation,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "service"),
			attribute.String("service.name", serviceName),
			attribute.String("operation", operation),
		)...),
	)
}

// StartRepoSpan starts a new span for a repository method.
func StartRepoSpan(ctx context.Context, tableName, operation string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "db."+tableName+"."+operation,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "repository"),
			attribute.String("db.system", "postgresql"),
			attribute.String("db.table", tableName),
			attribute.String("db.operation", operation),
		)...),
	)
}

// SpanFromContext returns the current span from context, if any.
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddEvent adds an event to the current span.
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	trace.SpanFromContext(ctx).AddEvent(name, trace.WithAttributes(attrs...))
}

// RecordError records an error on the current span.
func RecordError(ctx context.Context, err error, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	span.RecordError(err)
	span.SetAttributes(attrs...)
}

// SetAttributes sets attributes on the current span.
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	trace.SpanFromContext(ctx).SetAttributes(attrs...)
}

// TracingMiddleware is an HTTP middleware that adds tracing to requests.
// This is a convenience alias for otelhttp.NewHandler for use in router setup.
func TracingMiddleware(handler http.Handler) http.Handler {
	return otelhttp.NewHandler(handler, "http")
}
