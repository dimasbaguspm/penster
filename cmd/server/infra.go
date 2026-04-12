package main

import (
	"context"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/command"
	appquery "github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/internal/infrastructure/postgres"
	"github.com/dimasbaguspm/penster/internal/scheduler/engine"
	"github.com/dimasbaguspm/penster/pkg/observability"
)

type Infra struct {
	Scheduler *engine.Engine
	Server    *Server
}

func NewInfra(ctx context.Context, cfg *config.Config) *Infra {
	log := observability.NewLogger(ctx, "infra", "core")
	ctx, span := observability.StartServiceSpan(log.Context(), "core", "NewInfra")
	defer span.End()

	log.Info("infra starting", "auto_migrate", cfg.Migrate.AutoMigrate)

	log.Debug("infra connecting_db", "primary", cfg.DB.Primary, "max_conns", cfg.DB.MaxConns, "min_conns", cfg.DB.MinConns)
	conn := postgres.MustConnect(ctx, postgres.Config{
		Primary:  cfg.DB.Primary,
		MaxConns: cfg.DB.MaxConns,
		MinConns: cfg.DB.MinConns,
	})
	log.Debug("infra db_connected")

	if cfg.Migrate.AutoMigrate {
		log.Info("infra running_migration")
		postgres.RunMigration(ctx, postgres.Config{
			Primary: cfg.DB.Primary,
		})
		log.Info("infra migration_completed")
	}

	log.Debug("infra initializing_repositories")
	dbQueries := query.New(conn)

	accountRepo := repository.NewAccountRepository(dbQueries)
	categoryRepo := repository.NewCategoryRepository(dbQueries)
	rateCurrencyRepo := repository.NewRateCurrencyRepository(dbQueries)
	transactionRepo := repository.NewTransactionRepository(dbQueries)
	draftRepo := repository.NewDraftRepository(dbQueries)
	reportRepo := repository.NewReportRepository(dbQueries)
	log.Debug("infra repositories_initialized")

	log.Debug("infra initializing_queries")
	accountQuery := appquery.NewAccountQuery(accountRepo)
	accountCommand := command.NewAccountCommand(accountRepo)
	categoryQuery := appquery.NewCategoryQuery(categoryRepo)
	categoryCommand := command.NewCategoryCommand(categoryRepo)
	rateCurrencyQuery := appquery.NewRateCurrencyQuery(rateCurrencyRepo)
	rateCurrencyCommand := command.NewRateCurrencyCommand(rateCurrencyRepo)
	transactionQuery := appquery.NewTransactionQuery(transactionRepo)
	transactionCommand := command.NewTransactionCommand(transactionRepo)
	draftQuery := appquery.NewDraftQuery(draftRepo)
	draftCommand := command.NewDraftCommand(draftRepo)
	reportQuery := appquery.NewReportQuery(reportRepo)
	log.Debug("infra queries_initialized")

	log.Debug("infra initializing_services")
	accountService := service.NewAccountService(accountQuery, accountCommand)
	categoryService := service.NewCategoryService(categoryQuery, categoryCommand)
	rateCurrencyService := service.NewRateCurrencyService(rateCurrencyQuery, rateCurrencyCommand)
	transactionService := service.NewTransactionService(transactionQuery, transactionCommand, accountService, categoryService, rateCurrencyService, cfg)
	draftService := service.NewDraftService(draftQuery, draftCommand, accountService, categoryService, rateCurrencyService, transactionService, cfg)
	reportService := service.NewReportService(reportQuery)
	log.Debug("infra services_initialized")

	log.Debug("infra initializing_scheduler")
	scheduler := engine.NewEngine(cfg, rateCurrencyService)
	log.Debug("infra scheduler_initialized")

	log.Debug("infra initializing_server")
	server := NewServer(cfg, accountService, categoryService, transactionService, draftService, reportService)
	log.Debug("infra server_initialized")

	infra := &Infra{
		Scheduler: scheduler,
		Server:    server,
	}

	log.Info("infra started")
	return infra
}

func (i *Infra) Close(ctx context.Context) {
	log := observability.NewLogger(ctx, "infra", "core")
	ctx, span := observability.StartServiceSpan(log.Context(), "core", "Close")
	defer span.End()

	log.Info("close started")

	if i.Scheduler != nil {
		log.Debug("close stopping_scheduler")
		i.Scheduler.Stop()
		log.Debug("close scheduler_stopped")
	}

	if i.Server != nil {
		log.Debug("close stopping_server")
		i.Server.Stop(ctx)
		log.Debug("close server_stopped")
	}

	log.Info("close succeeded")
}
