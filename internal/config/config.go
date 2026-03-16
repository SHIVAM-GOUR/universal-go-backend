// Package config loads all application configuration from environment variables.
// It is the single source of truth for runtime configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all application configuration. Values are loaded exclusively
// from environment variables — no hardcoded defaults beyond safe dev fallbacks.
type Config struct {
	// HTTP server
	AppPort string
	AppEnv  string

	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBMaxConns int
	DBMinConns int
}

// Load reads environment variables and returns a validated Config.
// Returns an error if any required variable is missing or malformed.
func Load() (*Config, error) {
	dbUser, err := requireEnv("DB_USER")
	if err != nil {
		return nil, err
	}

	dbPassword, err := requireEnv("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	dbName, err := requireEnv("DB_NAME")
	if err != nil {
		return nil, err
	}

	maxConns, err := parseInt(getEnv("DB_MAX_CONNS", "25"), "DB_MAX_CONNS")
	if err != nil {
		return nil, err
	}

	minConns, err := parseInt(getEnv("DB_MIN_CONNS", "5"), "DB_MIN_CONNS")
	if err != nil {
		return nil, err
	}

	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
		DBMaxConns: maxConns,
		DBMinConns: minConns,
	}, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func requireEnv(key string) (string, error) {
	v := os.Getenv(key)
	if v == "" {
		return "", fmt.Errorf("required environment variable %q is not set", key)
	}
	return v, nil
}

func parseInt(s, name string) (int, error) {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q must be an integer: %w", name, err)
	}
	return v, nil
}
