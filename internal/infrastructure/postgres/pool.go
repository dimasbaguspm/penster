package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Primary  string
	MaxConns int
	MinConns int
}

func MustConnect(ctx context.Context, cfg Config) *pgxpool.Pool {
	log := observability.NewLogger(ctx, "infrastructure", "postgres")
	log.Info("Attempting to connect the database")

	config, err := pgxpool.ParseConfig(cfg.Primary)
	config.MinConns = int32(cfg.MinConns)
	config.MaxConns = int32(cfg.MaxConns)

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Error(fmt.Sprintf("Unable to connect with db, %v", err))
		os.Exit(1)
		return nil
	}

	log.Info("Connection established")
	return conn
}
