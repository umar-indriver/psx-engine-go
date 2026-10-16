package config

import (
	"os"
)

type Config struct {
	Port        string
	Environment string
	PSXBaseURL  string
	DatabaseURL string
}

func LoadConfig() *Config {
	port := getEnv("PORT", "8080")
	env := getEnv("ENVIRONMENT", "development")
	psxBaseURL := getEnv("PSX_BASE_URL", "https://dps.psx.com.pk")
	dbURL := getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/psx_db?sslmode=disable")

	return &Config{
		Port:        port,
		Environment: env,
		PSXBaseURL:  psxBaseURL,
		DatabaseURL: dbURL,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
