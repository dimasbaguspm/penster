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
	"github.com/dimasbaguspm/penster/internal/scheduler/jobs"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Infra struct {
	DB                  *pgxpool.Pool
	AccountService      *service.AccountService
	CategoryService     *service.CategoryService
	RateCurrencyService *service.RateCurrencyService
	Scheduler           *engine.Engine
}

func NewInfra(ctx context.Context, cfg *config.Config) (*Infra, error) {
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

	accountQuery := appquery.NewAccountQuery(accountRepo)
	accountCommand := command.NewAccountCommand(accountRepo)
	categoryQuery := appquery.NewCategoryQuery(categoryRepo)
	categoryCommand := command.NewCategoryCommand(categoryRepo)
	rateCurrencyQuery := appquery.NewRateCurrencyQuery(rateCurrencyRepo)
	rateCurrencyCommand := command.NewRateCurrencyCommand(rateCurrencyRepo)
	rateCurrencyService := service.NewRateCurrencyService(rateCurrencyQuery, rateCurrencyCommand)

	scheduler := engine.NewEngine(cfg, rateCurrencyService)

	return &Infra{
		DB:                  conn,
		AccountService:      service.NewAccountService(accountQuery, accountCommand),
		CategoryService:     service.NewCategoryService(categoryQuery, categoryCommand),
		RateCurrencyService: rateCurrencyService,
		Scheduler:           scheduler,
	}, nil
}

func (i *Infra) RegisterJobs(cfg *config.Config) {
	i.Scheduler.RegisterJob(jobs.NewRateCurrencyJob(cfg))
}

func (i *Infra) Close(ctx context.Context) {
	if i.Scheduler != nil {
		i.Scheduler.Stop()
	}
	if i.DB != nil {
		i.DB.Close()
	}
}

func (i *Infra) Health(ctx context.Context) error {
	return i.DB.Ping(ctx)
}
