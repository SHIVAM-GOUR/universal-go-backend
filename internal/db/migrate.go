package db

import (
	"embed"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"

	"my-app/internal/config"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// RunMigrations applies all pending up-migrations.
// It is only called when APP_ENV=production so local dev is never affected.
func RunMigrations(pool *pgxpool.Pool, cfg *config.Config) error {
	if cfg.AppEnv != "production" {
		return nil
	}

	sourceDriver, err := iofs.New(migrationFiles, "migrations")
	if err != nil {
		return fmt.Errorf("load migration files: %w", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	dbDriver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceDriver, cfg.DBName, dbDriver)
	if err != nil {
		return fmt.Errorf("create migrator: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("run migrations: %w", err)
	}

	return nil
}
