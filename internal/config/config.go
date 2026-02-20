// Package config provides configuration loading functionality.
package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Config holds the application configuration values.
type Config struct {
	Env        string
	DBUser    	string
	DBPassword 	string
	DBHost    	string
	DBPort		int
	DBName		string
	APIProtocol string
	APIHost 	string
	APIPort 	int
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
	env := os.Getenv("ENV")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPassword := os.Getenv("POSTGRES_PASSWORD")
	dbHost := os.Getenv("POSTGRES_HOST")
	dbPort := getIntEnv("POSTGRES_PORT", 5432)
	dbName := os.Getenv("POSTGRES_DB")
	apiProtocol := os.Getenv("API_PROTOCOL")
	apiHost := os.Getenv("API_HOST")
	apiPort := getIntEnv("API_PORT", 3000)
	if env != "dev" && env != "uat" && env != "prod" {
		logrus.Fatal("ENV environment variable must be set to 'dev', 'uat' or 'prod'")
	}
	if apiProtocol != "http" && apiProtocol != "https" {
		logrus.Fatal("API_PROTOCOL environment variable must be set to 'http' or 'https'")
	}
	if apiHost == "" {
		logrus.Fatal("API_PROTOCOL and API_HOST environment variables must be set")
	}
	return Config{
		Env: env,
		DBUser: dbUser,
		DBPassword: dbPassword,
		DBHost: dbHost,
		DBPort: dbPort,
		DBName: dbName,
		APIProtocol: apiProtocol,
		APIHost: apiHost,
		APIPort: apiPort,
	}
}
