package database

import (
	"context"
	"errors"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Threat struct {
	ID        int
	Filename  string
	Sha256    string
	Comments  []string
	Submitted string
}

type Comment struct {
	ID   int
	Text string
	Date string
}

func DatabaseConnect() (*pgx.Conn, error) {
	godotenv.Load(".env")

	database_string, exists := os.LookupEnv("DB_STRING")
	if !exists {
		return nil, errors.New("no db string configured within environment variables")
	}

	db, err := pgx.Connect(context.Background(), database_string)
	if err != nil {
		return nil, err
	}

	return db, nil
}
