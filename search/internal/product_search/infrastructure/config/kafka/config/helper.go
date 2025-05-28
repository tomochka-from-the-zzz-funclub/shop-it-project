package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func getEnv(key string) string {

	val := os.Getenv(key)

	if val == "" {

		log.Fatalf("[config:kafka][FATAL] missing environment variable %s", key)
	}

	log.Printf("[config:kafka] %s=%s", key, val)

	return val
}

func parseInt(key string) int {

	raw := getEnv(key)

	v, err := strconv.Atoi(raw)

	if err != nil {

		log.Fatalf("[config:kafka][FATAL] invalid int for %s: %v", key, err)
	}

	log.Printf("[config:kafka] parsed %s=%d", key, v)

	return v
}

func parseInt64(key string) int64 {

	raw := getEnv(key)

	v, err := strconv.ParseInt(raw, 10, 64)

	if err != nil {

		log.Fatalf("[config:kafka][FATAL] invalid int64 for %s: %v", key, err)

	}

	log.Printf("[config:kafka] parsed %s=%d", key, v)

	return v
}

func parseBool(key string) bool {

	raw := getEnv(key)

	v, err := strconv.ParseBool(raw)

	if err != nil {

		log.Fatalf("[config:kafka][FATAL] invalid bool for %s: %v", key, err)

	}

	log.Printf("[config:kafka] parsed %s=%t", key, v)

	return v
}

func parseDuration(key string) time.Duration {

	raw := getEnv(key)

	v, err := time.ParseDuration(raw)

	if err != nil {

		log.Fatalf("[config:kafka][FATAL] invalid duration for %s: %v", key, err)

	}

	log.Printf("[config:kafka] parsed %s=%s", key, v)

	return v
}

func parseHosts(key string) []string {

	raw := getEnv(key)

	hosts := strings.Split(raw, ",")

	log.Printf("[config:kafka] parsed %s=%v", key, hosts)

	return hosts
}
