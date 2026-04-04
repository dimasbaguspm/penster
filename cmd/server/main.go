package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimasbaguspm/penster/config"
	_ "github.com/dimasbaguspm/penster/docs"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	infra := NewInfra(ctx, cfg)
	defer infra.Close(ctx)

	infra.Scheduler.Start(ctx)
	infra.Server.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
