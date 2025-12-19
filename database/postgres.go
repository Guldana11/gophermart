package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func InitDB() *pgxpool.Pool {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DATABASE_URI")
	if dsn == "" {
		log.Fatal("DATABASE_URI is not set in environment variables")
	}

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	return pool
}
