// Package config provides configuration loading functionality.
package config

import (
	"log"
	"os"
)

// Config holds the application configuration values.
type Config struct {
	ApiPort string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() Config {
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		log.Fatal("API_PORT environment variable is required")
	}
	return Config{
		ApiPort: os.Getenv("API_PORT"),
	}
}
