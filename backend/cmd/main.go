package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/SpectreFury/odin-book/backend/internal/env"
)

func main() {
	err := env.Load()
	if err != nil {
		log.Fatal("ERROR: loading env", err)
		return
	}

	PORT := os.Getenv("PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Working"))
	})

	fmt.Println("Listening on PORT:", ":" + PORT)
	err = http.ListenAndServe(":" + PORT, mux)
	if err != nil {
		log.Fatal(err)
	}
}
