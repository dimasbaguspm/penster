package database

import (
	"context"
	"embed"
	"fmt"
	"log"

	"github.com/dimasbaguspm/penster/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Migrator handles database migrations
type Migrator struct {
	db  *pgx.Conn
	cfg *config.DBConfig
}

// NewMigrator creates a new Migrator
func NewMigrator(db *pgx.Conn, cfg *config.DBConfig) *Migrator {
	return &Migrator{db: db, cfg: cfg}
}

// MigrateUp runs all pending migrations
func (m *Migrator) MigrateUp(ctx context.Context) error {
	log.Println("Running database migrations...")

	// Get migration files from the embedded filesystem
	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	// Create migrator using pgx driver
	migrator, err := migrate.NewWithSourceInstance("iofs", source, m.cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	// Run migrations
	if err := migrator.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}

// MigrateDown rolls back all migrations
func (m *Migrator) MigrateDown(ctx context.Context) error {
	source, err := iofs.New(migrationFS, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source: %w", err)
	}

	migrator, err := migrate.NewWithSourceInstance("iofs", source, m.cfg.DSN())
	if err != nil {
		return fmt.Errorf("failed to create migrator: %w", err)
	}
	defer migrator.Close()

	if err := migrator.Down(); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}

	return nil
}
