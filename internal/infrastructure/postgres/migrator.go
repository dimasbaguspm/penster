package postgres

import (
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunMigration(cfg Config) {
	slog.Info("[Migrator]: trying to migrate tables into DB")

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	if migrationsPath == "" {
		execPath, err := os.Executable()
		if err != nil {
			slog.Error("[Migrator]: failed to get executable path", "error", err)
			os.Exit(1)
		}
		// Resolve symlinks (e.g., from Air live-reload)
		execPath, err = filepath.EvalSymlinks(execPath)
		if err != nil {
			slog.Error("[Migrator]: failed to resolve symlinks", "error", err)
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
			slog.Error("[Migrator]: migrations folder not found")
			os.Exit(1)
		}
	}

	migrationsURL := "file://" + migrationsPath
	m, err := migrate.New(migrationsURL, cfg.Primary)
	if err != nil {
		slog.Error("[Migrator]: migration failed something odd while lookup the migrations file", "error", err)
		os.Exit(1)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("[Migrator]: migration success without no change!")
			return
		}

		slog.Error("[Migrator]: unable to migrate the db", "error", err)
		os.Exit(1)
	}

	slog.Info("[Migrator]: success to migrate the latest version!")
}
