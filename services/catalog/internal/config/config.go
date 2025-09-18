package config

import (
	"os"
)

type Config struct {
	HTTPPort       string
	DBURL          string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
}

func Load() *Config {
	cfg := &Config{}

	cfg.HTTPPort = getEnv("CATALOG_PORT", "8082")
	cfg.DBURL = getEnv("DATABASE_URL", "postgres://app:app@localhost:5432/app?sslmode=disable")
	cfg.MinioEndpoint = getEnv("MINIO_ENDPOINT", "minio:9000")
	cfg.MinioAccessKey = getEnv("MINIO_ACCESS_KEY", "minioadmin")
	cfg.MinioSecretKey = getEnv("MINIO_SECRET_KEY", "minioadmin123")
	cfg.MinioBucket = getEnv("MINIO_BUCKET", "product-images")

	return cfg
}

func getEnv(key string, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
