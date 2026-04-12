package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/interface/handler"
	"github.com/dimasbaguspm/penster/internal/interface/router"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type Server struct {
	srv *http.Server
}

func NewServer(cfg *config.Config, accountSvc *service.AccountService, categorySvc *service.CategoryService, transactionSvc *service.TransactionService, draftSvc *service.DraftService, reportSvc *service.ReportService) *Server {
	healthHandler := handler.NewHealthHandler(nil)
	r := router.NewRouter(healthHandler, accountSvc, categorySvc, transactionSvc, draftSvc, reportSvc)

	return &Server{
		srv: &http.Server{
			Addr:    fmt.Sprintf(":%s", cfg.App.Port),
			Handler: r.Routes(),
		},
	}
}

func (s *Server) Start(ctx context.Context) {
	log := observability.NewLogger(ctx, "core", "server")
	log.Info("Starting server", "addr", s.srv.Addr)
	go func() {
		if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
		}
	}()
}

func (s *Server) Stop(ctx context.Context) error {
	log := observability.NewLogger(ctx, "core", "server")
	log.Info("Shuting down")
	return s.srv.Shutdown(ctx)
}
