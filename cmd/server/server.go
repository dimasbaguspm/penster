package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/handler"
	"github.com/dimasbaguspm/penster/internal/interface/router"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *config.Config, accountSvc *service.AccountService, categorySvc *service.CategoryService, transactionSvc *service.TransactionService) *Server {
	healthHandler := handler.NewHealthHandler(nil)
	r := router.NewRouter(healthHandler, accountSvc, categorySvc, transactionSvc)

	return &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.App.Port),
			Handler: r.Routes(),
		},
	}
}

func (s *Server) Start() {
	go func() {
		log.Printf("Starting server on :%s", s.srv.Addr)
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
