package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Init() *pgxpool.Pool {
	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		dsn = "postgres://postgres:postgres@localhost:5432/loyalty?sslmode=disable"
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	DB = pool
	return pool
}
