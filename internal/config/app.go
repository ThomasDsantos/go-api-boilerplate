package config

import (
	"os"

	"github.com/rs/zerolog"
)

type ServerConfig struct {
	Port        string
	LogLevel    zerolog.Level
	Environment string
	ServiceName string
	APIBasePath string
}

func Load() (*ServerConfig, error) {
	port := getEnv("SERVER_PORT", "8080")
	logLevelStr := getEnv("LOG_LEVEL", "info")
	environment := getEnv("ENVIRONMENT", "local")
	apiBasePath := getEnv("API_BASE_PATH", "/v1")
	serviceName := getEnv("SERVICE_NAME", "api")

	logLevel, err := zerolog.ParseLevel(logLevelStr)
	if err != nil {
		// Default to info level if parsing fails, or handle error more strictly
		logLevel = zerolog.InfoLevel
	}

	return &ServerConfig{
		Port:        port,
		LogLevel:    logLevel,
		Environment: environment,
		APIBasePath: apiBasePath,
		ServiceName: serviceName,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
