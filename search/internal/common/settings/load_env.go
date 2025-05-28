package settings

import (
	"github.com/joho/godotenv"
	"log"
)

func LoadEnv(path string) error {
	log.Printf("Loading environment from %s", path)
	if err := godotenv.Load(path); err != nil {
		log.Printf("ERROR: failed to load env file %s: %v", path, err)
		return err
	}
	log.Println("Environment variables loaded successfully")
	return nil
}
