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

	cwd, err := os.Getwd()
	if err != nil {
		slog.Error("[Migrator]: failed to get current working directory", "error", err)
		os.Exit(1)
	}
	migrationsPath, err := filepath.Abs(filepath.Join(cwd, "migrations"))
	if err != nil {
		slog.Error("[Migrator]: failed to resolve absolute path to migrations", "error", err)
		os.Exit(1)
	}

	m, err := migrate.New("file:///"+migrationsPath, cfg.Primary)
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
