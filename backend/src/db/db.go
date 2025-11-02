package db

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

func  Connect() (*pgx.Conn, error) {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URI"))
	if err != nil {
		return nil, err
	}

	log.Print("Database connected")
	return conn, nil
}
