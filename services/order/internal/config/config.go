package config

import (
	"os"
)

type Config struct {
	HTTPPort    string
	DBURL       string
	CatalogURL  string
	CatalogPort string
}

func Load() *Config {
	cfg := &Config{}

	cfg.HTTPPort = getEnv("ORDER_PORT", "8083")
	cfg.DBURL = getEnv("DATABASE_URL", "postgres://app:app@localhost:5432/app?sslmode=disable")
	cfg.CatalogURL = getEnv("CATALOG_URL", "localhost")
	cfg.CatalogPort = getEnv("CATALOG_PORT", "8082")

	return cfg
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
