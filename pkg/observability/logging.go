package observability

import (
	"context"
	"log/slog"
	"os"

	"github.com/dimasbaguspm/penster/config"
	"github.com/google/uuid"
)

var logger *slog.Logger

type contextKey string

const txnIDKey contextKey = "txn_id"

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

func SlogLogger() *slog.Logger {
	return logger
}

func GenTransactionID(ctx context.Context) context.Context {
	return context.WithValue(ctx, txnIDKey, uuid.New().String())
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
