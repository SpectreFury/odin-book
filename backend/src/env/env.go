package env

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	log.Print("Loaded env")
	return nil
}
