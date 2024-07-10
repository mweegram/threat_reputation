package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func DatabaseConnect() *pgxpool.Pool {
	godotenv.Load(".env")

	database_string, exists := os.LookupEnv("DB_STRING")
	if !exists {
		log.Printf("Error: %s", "no db string configured within environment variables")
		return nil
	}

	db, err := pgxpool.New(context.Background(), database_string)
	if err != nil {
		log.Printf("Error: %v", err)
		return nil
	}

	return db
}
