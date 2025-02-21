package config

import (
	"os"
	"path/filepath"
)

type Config struct {
	DBPath      string
	ServerPort  string
	Environment string
}

func New() *Config {
	return &Config{
		DBPath:      getEnvOrDefault("DB_PATH", filepath.Join(".", "words.db")),
		ServerPort:  getEnvOrDefault("SERVER_PORT", "8080"),
		Environment: getEnvOrDefault("ENV", "development"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
} 