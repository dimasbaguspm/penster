package observability

import (
	"context"
	"log/slog"
	"os"

	"github.com/dimasbaguspm/penster/config"
	"github.com/google/uuid"
)

var logger *slog.Logger
var observabilityEnabled bool

type contextKey string

const txnIDKey contextKey = "txn_id"

func InitLogger(ctx context.Context, cfg *config.Config) *slog.Logger {
	observabilityEnabled = cfg.Observability.Enabled

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger = slog.New(handler)
	slog.SetDefault(logger)

	if observabilityEnabled {
		slog.Info("logger initialized",
			"env", cfg.App.Env,
		)
	} else {
		slog.Info("logger disabled",
			"env", cfg.App.Env,
		)
	}

	return logger
}

func GetTransactionID(ctx context.Context) string {
	if v := ctx.Value(txnIDKey); v != nil {
		return v.(string)
	}
	return ""
}

type Logger struct {
	ctx       context.Context
	layer     string
	component string
	log       *slog.Logger
}

func NewLogger(ctx context.Context, layer, component string) *Logger {
	if GetTransactionID(ctx) == "" {
		ctx = context.WithValue(ctx, txnIDKey, uuid.New().String())
	}
	return &Logger{
		ctx:       ctx,
		layer:     layer,
		component: component,
		log:       logger,
	}
}

func (l *Logger) Info(msg string, attrs ...any) {
	l.Log(slog.LevelInfo, msg, attrs...)
}

func (l *Logger) Error(msg string, attrs ...any) {
	l.Log(slog.LevelError, msg, attrs...)
}

func (l *Logger) Warn(msg string, attrs ...any) {
	l.Log(slog.LevelWarn, msg, attrs...)
}

func (l *Logger) Debug(msg string, attrs ...any) {
	l.Log(slog.LevelDebug, msg, attrs...)
}

func (l *Logger) Log(level slog.Level, msg string, attrs ...any) {
	args := []any{"layer", l.layer, "component", l.component, "txn_id", GetTransactionID(l.ctx)}
	args = append(args, attrs...)
	l.log.Log(l.ctx, level, msg, args...)
}

func (l *Logger) WithCtx(ctx context.Context) *Logger {
	return &Logger{
		ctx:       ctx,
		layer:     l.layer,
		component: l.component,
		log:       l.log,
	}
}

func (l *Logger) Context() context.Context {
	return l.ctx
}
