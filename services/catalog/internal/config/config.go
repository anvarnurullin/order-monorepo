package config

import (
	"os"
)

type Config struct {
	HTTPPort    string
	DBURL       string
}

func Load() *Config {
	cfg := &Config{}

	cfg.HTTPPort = getEnv("CATALOG_HTTP_PORT", "8082")
	cfg.DBURL = getEnv("DATABASE_URL", "postgres://app:app@localhost:5433/app?sslmode=disable")

	return cfg
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
