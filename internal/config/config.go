// Package config provides configuration loading functionality.
package config

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// Config holds the application configuration values.
type Config struct {
	Env        string `env:"ENV,required,allowed:dev|uat|prod"`
	DBUser    	string `env:"POSTGRES_USER,required"`
	DBPassword 	string `env:"POSTGRES_PASSWORD,required"`
	DBHost    	string `env:"POSTGRES_HOST,required"`
	DBPort		int `env:"POSTGRES_PORT,required"`
	DBName		string `env:"POSTGRES_DB,required"`
	APIProtocol string `env:"API_PROTOCOL,required,allowed:http|https"`
	APIHost 	string `env:"API_HOST,required"`
	APIPort 	int `env:"API_PORT,required"`
	APIDeliveryIntervalSec int `env:"API_DELIVERY_INTERVAL_SEC,default:30"`
	APIRecoveryIntervalSec int `env:"API_RECOVERY_INTERVAL_SEC,default:300"`
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
	var config Config
	err := loadEnvs(&config)
	if err != nil {
		logrus.WithError(err).Fatal("Error loading configuration")
	}
	fmt.Println(config)

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
	apiDeliveryIntervalSec := getIntEnv("API_DELIVERY_INTERVAL_SEC", 30)
	apiRecoveryIntervalSec := getIntEnv("API_RECOVERY_INTERVAL_SEC", 300)
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
		APIDeliveryIntervalSec: apiDeliveryIntervalSec,
		APIRecoveryIntervalSec: apiRecoveryIntervalSec,
	}
}
