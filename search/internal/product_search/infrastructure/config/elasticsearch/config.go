package elasticsearch

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Address             string // host:port, например "127.0.0.1:9200"
	Username            string // имя пользователя, например "elastic"
	Password            string
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
}

func NewConfig(
	address, username, password string,
	maxIdle, maxPerHost int,
	idleTimeout time.Duration,
) *Config {
	return &Config{
		Address:             address,
		Username:            username,
		Password:            password,
		MaxIdleConns:        maxIdle,
		MaxIdleConnsPerHost: maxPerHost,
		IdleConnTimeout:     idleTimeout,
	}
}

func ReadConfigFromEnv() *Config {
	log.Println("[elasticsearch] reading config from env")

	address := getEnv("ES_PORT")
	username := getEnv("ELASTIC_USER")
	password := getEnv("ELASTIC_PASSWORD")

	maxIdle := parseEnvInt("ES_MAX_IDLE_CONNS")
	maxPerHost := parseEnvInt("ES_MAX_IDLE_CONNS_PER_HOST")
	idleSec := parseEnvInt("ES_IDLE_CONN_TIMEOUT")

	cfg := NewConfig(
		address,
		username,
		password,
		maxIdle,
		maxPerHost,
		time.Duration(idleSec)*time.Second,
	)

	log.Printf(
		"[elasticsearch] config loaded: address=%s product=%s maxIdle=%d maxPerHost=%d idleTimeout=%s",
		cfg.Address, cfg.Username,
		cfg.MaxIdleConns, cfg.MaxIdleConnsPerHost,
		cfg.IdleConnTimeout,
	)

	return cfg
}

func getEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Panicf("[elasticsearch][FATAL] env %s not set", key)
	}
	return v
}

func parseEnvInt(key string) int {
	v := getEnv(key)
	i, err := strconv.Atoi(v)
	if err != nil {
		log.Panicf("[elasticsearch][FATAL] invalid %s: %v", key, err)
	}
	return i
}
