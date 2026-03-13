package main

import (
	"context"
	"net/http"
	"os"

	"github.com/SpectreFury/odin-book/backend/cmd/handler"
	"github.com/SpectreFury/odin-book/backend/internal/env"
	"github.com/SpectreFury/odin-book/backend/internal/migration"
	"github.com/SpectreFury/odin-book/backend/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log := logger.Logger{}

	err := env.Load()
	if err != nil {
		log.Error("Unable to load env")
		return
	}

	PORT := os.Getenv("PORT")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		log.Error("ERROR: connecting to database")
		return
	}
	defer conn.Close()
	log.Log("Connected to database")

	migration.RunMigration(conn, "migrations")

	authHandler := handler.AuthHandler{
		DB: conn,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", authHandler.SignupHandler)

	log.Log("Listening on PORT 4000")
	err = http.ListenAndServe(":"+PORT, mux)
	if err != nil {
		log.Error(err.Error())
	}

}
