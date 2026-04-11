package observability

import (
	"context"
	"net/http"
	"strings"

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

// StartCommandSpan starts a new span for a command method.
func StartCommandSpan(ctx context.Context, entity, operation string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, entity+"Command."+operation,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "command"),
			attribute.String("entity", entity),
			attribute.String("operation", strings.ToLower(operation)),
		)...),
	)
}

// StartQuerySpan starts a new span for a query method.
func StartQuerySpan(ctx context.Context, entity, operation string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, entity+"Query."+operation,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "query"),
			attribute.String("entity", entity),
			attribute.String("operation", strings.ToLower(operation)),
		)...),
	)
}

// StartValueObjectSpan starts a new span for a valueobject helper.
func StartValueObjectSpan(ctx context.Context, entity, function string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "valueobjects."+entity+"."+function,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "valueobject"),
			attribute.String("entity", entity),
			attribute.String("function", function),
		)...),
	)
}

// StartHandlerSpan starts a new span for an HTTP handler method.
func StartHandlerSpan(ctx context.Context, handler, operation string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, handler+"Handler."+operation,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "handler"),
			attribute.String("handler.name", handler),
			attribute.String("operation", strings.ToLower(operation)),
		)...),
	)
}

// StartJobSpan starts a new span for a scheduler job.
func StartJobSpan(ctx context.Context, jobName string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "job."+jobName,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "job"),
			attribute.String("job.name", jobName),
		)...),
	)
}

// StartDTOSpan starts a new span for a DTO parsing or validation function.
func StartDTOSpan(ctx context.Context, entity, function string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	return Tracer().Start(ctx, "dto."+entity+"."+function,
		trace.WithAttributes(append(attrs,
			attribute.String("layer", "dto"),
			attribute.String("entity", entity),
			attribute.String("function", function),
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
