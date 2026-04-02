package main

import (
	"context"
	"fmt"
	"log"

	"github.com/dimasbaguspm/penster/config"
	"github.com/dimasbaguspm/penster/internal/application/command"
	appquery "github.com/dimasbaguspm/penster/internal/application/query"
	"github.com/dimasbaguspm/penster/internal/application/service"
	"github.com/dimasbaguspm/penster/internal/domain/repository"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database"
	"github.com/dimasbaguspm/penster/internal/infrastructure/database/query"
	"github.com/jackc/pgx/v5"
)

type Infra struct {
	DB              *pgx.Conn
	AccountService  *service.AccountService
	CategoryService *service.CategoryService
}

func NewInfra(ctx context.Context, cfg *config.Config) (*Infra, error) {
	conn, err := pgx.Connect(ctx, cfg.DB.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Connected to database")

	if cfg.Migrate.AutoMigrate {
		migrator := database.NewMigrator(conn, &cfg.DB)
		if err := migrator.MigrateUp(ctx); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	dbQueries := query.New(conn)

	accountRepo := repository.NewAccountRepository(dbQueries)
	categoryRepo := repository.NewCategoryRepository(dbQueries)

	accountQuery := appquery.NewAccountQuery(accountRepo)
	accountCommand := command.NewAccountCommand(accountRepo)
	categoryQuery := appquery.NewCategoryQuery(categoryRepo)
	categoryCommand := command.NewCategoryCommand(categoryRepo)

	return &Infra{
		DB:              conn,
		AccountService:  service.NewAccountService(accountQuery, accountCommand),
		CategoryService: service.NewCategoryService(categoryQuery, categoryCommand),
	}, nil
}

func (i *Infra) Close(ctx context.Context) {
	if i.DB != nil {
		i.DB.Close(ctx)
	}
}

func (i *Infra) Health(ctx context.Context) error {
	return i.DB.Ping(ctx)
}
