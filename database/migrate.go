package database

import (
	"database/sql"
	"fmt"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func Migrate(dbURL string) error {
	sqlDB, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migrate driver: %w", err)
	}

	relPath := "./migrations"

	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute migrations path: %w", err)
	}

	sourceURL := "file://" + absPath

	fmt.Println("migration source:", sourceURL)

	m, err := migrate.NewWithDatabaseInstance(
		sourceURL,
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}
