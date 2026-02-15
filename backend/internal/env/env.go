package env

import (
	"github.com/joho/godotenv"
)

func Load() error {
	// Load env
	err := godotenv.Load()
	if err != nil {
		return err
	}

	return nil
}
