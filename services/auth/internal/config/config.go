package config

import (
	"os"
)

type Config struct {
	HTTPPort  string
	DBURL     string
	JWTSecret string
}

func Load() *Config {
	cfg := &Config{}

	cfg.HTTPPort = getEnv("AUTH_PORT", "8084")
	cfg.DBURL = getEnv("DATABASE_URL", "postgres://app:app@localhost:5432/app?sslmode=disable")
	cfg.JWTSecret = getEnv("JWT_SECRET", "your-secret-key-change-in-production")

	return cfg
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}