package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser     string
	DBPassword string
	JwtSecret  string
	DBName     string
	DBHost     string
	DBPort     string
	SslMode    string
}

func LoadConfig() Config {
	if os.Getenv("ENV") != "docker" {
		if err := godotenv.Load("configs/local.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	return Config{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		JwtSecret:  os.Getenv("JWT_SECRET"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		SslMode:    os.Getenv("DB_SSLMODE"),
	}
}
