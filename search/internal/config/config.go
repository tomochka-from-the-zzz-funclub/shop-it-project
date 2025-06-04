package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser             string
	DBPassword         string
	DBName             string
	DBHost             string
	DBPort             string
	SslMode            string
	KafkaHost          string
	KafkaPort          string
	KafkaTopic         string
	KafkaResponseTopic string // Новое поле для топика ответов
	TTL                time.Duration
}

func LoadConfig() Config {
	if os.Getenv("ENV") != "docker" {
		if err := godotenv.Load("configs/local.env"); err != nil {
			log.Fatal("Error loading .env file")
		}
	}
	TTLstr := os.Getenv("TTL")
	TTL, err := time.ParseDuration(TTLstr)
	if err != nil {
		log.Fatal("Error during parsing ttl: " + err.Error())
	}
	return Config{
		DBUser:             os.Getenv("DB_USER"),
		DBPassword:         os.Getenv("DB_PASSWORD"),
		DBName:             os.Getenv("DB_NAME"),
		DBHost:             os.Getenv("DB_HOST"),
		DBPort:             os.Getenv("DB_PORT"),
		SslMode:            os.Getenv("DB_SSLMODE"),
		KafkaPort:          os.Getenv("KAFKA_PORT"),
		KafkaHost:          os.Getenv("KAFKA_HOST"),
		KafkaTopic:         os.Getenv("KAFKA_TOPIC"),
		KafkaResponseTopic: os.Getenv("KAFKA_RESPONSE_TOPIC"), // Считывание нового топика
		TTL:                TTL,
	}
}
