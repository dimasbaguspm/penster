package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimasbaguspm/penster/config"
	_ "github.com/dimasbaguspm/penster/docs"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()
	obs := observability.Init(ctx, cfg)
	defer obs.Shutdown(ctx)

	infra := NewInfra(ctx, cfg)
	defer infra.Close(ctx)

	infra.Scheduler.Start(ctx)
	infra.Server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server")
}