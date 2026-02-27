// Package config provides configuration loading functionality.
package config

import (
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
	APIDeliveryIntervalSec int `env:"API_DELIVERY_INTERVAL_SEC,default:30,min:10,max:60"`
	APIRecoveryIntervalSec int `env:"API_RECOVERY_INTERVAL_SEC,default:300,mim:300,max:3600"`
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() Config {
	var config Config
	err := loadEnvs(&config)
	if err != nil {
		logrus.WithError(err).Fatal("Error loading configuration")
	}
	return config
}
