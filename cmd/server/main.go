package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dimasbaguspm/penster/config"
	_ "github.com/dimasbaguspm/penster/docs"
	"github.com/dimasbaguspm/penster/internal/interface/handler"
	"github.com/dimasbaguspm/penster/internal/interface/router"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.Load()

	infra, err := NewInfra(ctx, cfg)
	if err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}
	defer infra.Close(ctx)

	infra.RegisterJobs(cfg)
	infra.Scheduler.Start(ctx)

	healthHandler := handler.NewHealthHandler(infra)
	r := router.NewRouter(healthHandler, infra.AccountService, infra.CategoryService)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.App.Port),
		Handler: r.Routes(),
	}

	go func() {
		log.Printf("Starting server on :%s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
}
