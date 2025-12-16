package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Migrate(db *pgxpool.Pool) {
	migrationsPath := "./migrations"

	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		log.Fatalf("failed to read migrations folder: %v", err)
	}

	ctx := context.Background()

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".sql" {
			continue
		}
		path := filepath.Join(migrationsPath, file.Name())
		sqlBytes, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("failed to read migration file %s: %v", path, err)
		}

		_, err = db.Exec(ctx, string(sqlBytes))
		if err != nil {
			log.Fatalf("failed to execute migration %s: %v", path, err)
		}

		fmt.Printf("migration %s applied\n", file.Name())
	}
}
