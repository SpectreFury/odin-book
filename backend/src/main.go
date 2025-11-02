package main

import (
	"log"

	"github.com/SpectreFury/odin-book/backend/src/db"
	"github.com/SpectreFury/odin-book/backend/src/env"
	"github.com/SpectreFury/odin-book/backend/src/router"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatalf("Error loading env: %v", err)
	}

	conn, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}


	log.Print("Connection: ", conn)
	r := router.Init()
	router.RegisterRoutes(r, conn)
	router.Run(r)
}
