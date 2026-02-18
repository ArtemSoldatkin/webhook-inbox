// Package config provides configuration loading functionality.
package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Config holds the application configuration values.
type Config struct {
	DBUser    	string
	DBPassword 	string
	DBHost    	string
	DBPort		int
	DBName		string
	ApiPort 	string
}

// getIntEnv retrieves an integer environment variable or returns a default value.
func getIntEnv(envVar string, defaultValue int) int {
	valueStr := os.Getenv(envVar)
	if valueStr == "" {
		return defaultValue
	}
	var value int
	_, err := fmt.Sscanf(valueStr, "%d", &value)
	if err != nil {
		logrus.WithError(err).Fatalf("Invalid value for %s", envVar)
	}
	return value
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() Config {
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := getIntEnv("POSTGRES_PORT", 5432)
	dbName := os.Getenv("POSTGRES_DB")
	apiPort := os.Getenv("API_PORT")
	if apiPort == "" {
		logrus.Fatal("API_PORT environment variable is required")
	}
	return Config{
		DBUser:    dbUser,
		DBPassword: dbPassword,
		DBHost:     dbHost,
		DBPort:    dbPort,
		DBName:    dbName,
		ApiPort: os.Getenv("API_PORT"),
	}
}
