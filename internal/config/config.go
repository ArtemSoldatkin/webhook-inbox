// Package config provides configuration loading functionality.
package config

import (
	"os"
)

// Config holds the application configuration values.
type Config struct {
	API_PORT string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() Config {
	return Config{
		API_PORT: os.Getenv("API_PORT"),
	}
}
