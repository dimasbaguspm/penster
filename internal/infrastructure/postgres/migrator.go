package postgres

import (
	"context"
	"errors"
	"os"
	"path/filepath"

	"github.com/dimasbaguspm/penster/pkg/observability"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(ctx context.Context, cfg Config) {
	log := observability.NewLogger(ctx, "infrastructure", "migrator")
	log.Info("trying to migrate tables into DB")

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		execPath, err := os.Executable()
		if err != nil {
			log.Error("failed to get executable path", "error", err)
			os.Exit(1)
		}
		// Resolve symlinks (e.g., from Air live-reload)
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			log.Error("failed to resolve symlinks", "error", err)
			os.Exit(1)
		}
		execDir := filepath.Dir(execPath)

		// Try executable's dir first, then fallback to /app/migrations
		for _, mp := range []string{
			filepath.Join(execDir, "migrations"),
			"/app/migrations",
		} {
			if _, err := os.Stat(mp); err == nil {
				migrationsPath = mp
				break
			}
		}
		if migrationsPath == "" {
			log.Error("migrations folder not found")
			os.Exit(1)
		}
	}

	migrationsURL := "file://" + migrationsPath
	m, err := migrate.New(migrationsURL, cfg.Primary)
	if err != nil {
		log.Error("migration failed something odd while lookup the migrations file", "error", err)
		os.Exit(1)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("migration success without no change!")
			return
		}

		log.Error("unable to migrate the db", "error", err)
		os.Exit(1)
	}

	log.Info("success to migrate the latest version!")
}
