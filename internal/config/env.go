package config

import (
	"os"

	"currency-exchange/pkg/postgres"
)

func GetEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func HTTPAddr() string {
	return GetEnvOrDefault("HTTP_ADDR", ":8080")
}

func PostgresConfigFromEnv() postgres.Config {
	return postgres.Config{
		Host:     GetEnvOrDefault("PG_HOST", "localhost"),
		Port:     GetEnvOrDefault("PG_PORT", "5432"),
		User:     GetEnvOrDefault("PG_USER", "postgres"),
		Password: GetEnvOrDefault("PG_PASSWORD", "postgres"),
		DBName:   GetEnvOrDefault("PG_DBNAME", "postgres"),
		SSLMode:  GetEnvOrDefault("PG_SSLMODE", "disable"),
	}
}
