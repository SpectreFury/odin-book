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
	err := env.Load()
	if err != nil {
		logger.Error("Unable to load env")
		return
	}

	PORT := os.Getenv("PORT")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := pgxpool.New(context.Background(), DATABASE_URL)
	if err != nil {
		logger.Error("ERROR: connecting to database")
		return
	}
	defer conn.Close()
	logger.Info("Connected to database")

	migration.RunMigration(conn, "migrations")

	authHandler := handler.AuthHandler{
		DB: conn,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", authHandler.SignupHandler)
	mux.HandleFunc("POST /login", authHandler.LoginHandler)

	logger.Info("Listening on PORT 4000")
	err = http.ListenAndServe(":"+PORT, mux)
	if err != nil {
		logger.Error(err.Error())
	}

}
