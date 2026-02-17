package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SpectreFury/odin-book/backend/cmd/handler"
	"github.com/SpectreFury/odin-book/backend/internal/env"
	"github.com/SpectreFury/odin-book/backend/internal/migration"
	"github.com/jackc/pgx/v5"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal(err)
		return
	}

	PORT := os.Getenv("PORT")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	conn, err := pgx.Connect(context.Background(), DATABASE_URL)
	if err != nil {
		log.Fatal("ERROR: connecting to database")
		return
	}
	defer conn.Close(context.Background())
	fmt.Println("Connected to database")

	migration.RunMigration(conn, "migrations")

	authHandler := handler.AuthHandler{}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /signup", authHandler.SignupHandler)

	fmt.Println("Listening on PORT:", PORT)
	err = http.ListenAndServe(":"+PORT, mux)
	if err != nil {
		log.Fatal(err)
	}

}
