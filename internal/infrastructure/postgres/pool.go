package postgres

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/dimasbaguspm/penster/pkg/observability"
)

type Config struct {
	Primary  string
	MaxConns int
	MinConns int
}

// Pool wraps pgxpool.Pool and implements PoolStats interface.
type Pool struct {
	*pgxpool.Pool
}

// Stat returns pool statistics for observability gauges.
func (p *Pool) Stat(ctx context.Context) observability.PoolsStat {
	stat := p.Pool.Stat()
	return observability.PoolsStat{
		Acquired: int64(stat.AcquiredConns()),
		Idle:     int64(stat.IdleConns()),
		Total:    int64(stat.TotalConns()),
	}
}

func MustConnect(ctx context.Context, cfg Config) *Pool {
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
	return &Pool{Pool: conn}
}
