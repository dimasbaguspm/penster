package main

import (
	"context"
	"log"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/command"
	appquery "github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/dimasbaguspm/penster/internal/infrastructure/postgres"
	"github.com/dimasbaguspm/penster/internal/scheduler/engine"
)

type Infra struct {
	Scheduler *engine.Engine
	Server    *Server
}

func NewInfra(ctx context.Context, cfg *config.Config) *Infra {
	conn := postgres.MustConnect(ctx, postgres.Config{
		Primary:  cfg.DB.Primary,
		MaxConns: cfg.DB.MaxConns,
		MinConns: cfg.DB.MinConns,
	})

	log.Println("Connected to database")

	if cfg.Migrate.AutoMigrate {
		postgres.RunMigration(postgres.Config{
			Primary: cfg.DB.Primary,
		})
	}

	dbQueries := query.New(conn)

	accountRepo := repository.NewAccountRepository(dbQueries)
	categoryRepo := repository.NewCategoryRepository(dbQueries)
	rateCurrencyRepo := repository.NewRateCurrencyRepository(dbQueries)
	transactionRepo := repository.NewTransactionRepository(dbQueries)
	draftRepo := repository.NewDraftRepository(dbQueries)
	reportRepo := repository.NewReportRepository(dbQueries)

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

	accountService := service.NewAccountService(accountQuery, accountCommand)
	categoryService := service.NewCategoryService(categoryQuery, categoryCommand)
	rateCurrencyService := service.NewRateCurrencyService(rateCurrencyQuery, rateCurrencyCommand)
	transactionService := service.NewTransactionService(transactionQuery, transactionCommand, accountService, categoryService, rateCurrencyService, cfg)
	draftService := service.NewDraftService(draftQuery, draftCommand, accountService, categoryService, rateCurrencyService, transactionService, cfg)
	reportService := service.NewReportService(reportQuery)

	scheduler := engine.NewEngine(cfg, rateCurrencyService)
	server := NewServer(cfg, accountService, categoryService, transactionService, draftService, reportService)

	return &Infra{
		Scheduler: scheduler,
		Server:    server,
	}
}

func (i *Infra) Close(ctx context.Context) {
	if i.Scheduler != nil {
		i.Scheduler.Stop()
	}
	if i.Server != nil {
		i.Server.Stop(ctx)
	}
}
