// Package config provides configuration loading functionality.
package config

import (
	"github.com/sirupsen/logrus"
)

// Config holds the application configuration values.
type Config struct {
	Env string `env:"ENV,required,allowed:dev|uat|prod"`

	DBUser     string `env:"POSTGRES_USER,required"`
	DBPassword string `env:"POSTGRES_PASSWORD,required"`
	DBHost     string `env:"POSTGRES_HOST,required"`
	DBPort     int    `env:"POSTGRES_PORT,required"`
	DBName     string `env:"POSTGRES_DB,required"`

	APIProtocol string `env:"API_PROTOCOL,required,allowed:http|https"`
	APIHost     string `env:"API_HOST,required"`
	APIPort     int    `env:"API_PORT,required"`

	APIDeliveryIntervalSec       int `env:"API_DELIVERY_INTERVAL_SEC,default:30,min:10,max:60"`
	APIDeliveryMaxConcurrency    int `env:"API_DELIVERY_MAX_CONCURRENCY,default:10,min:1,max:100"`
	APIDeliveryTimeoutSec        int `env:"API_DELIVERY_TIMEOUT_SEC,default:15,min:5,max:60"`
	APIDeliveryRequestTimeoutSec int `env:"API_DELIVERY_REQUEST_TIMEOUT_SEC,default:15,min:5,max:60"`

	APIDeliveryMaxRetries          int `env:"API_DELIVERY_MAX_RETRIES,default:3,min:0,max:10"`
	APIDeliveryRetryBackoffBaseSec int `env:"API_DELIVERY_RETRY_BACKOFF_BASE_SEC,default:1,min:1,max:60"`
	APIDeliveryRetryBackoffMaxSec  int `env:"API_DELIVERY_RETRY_BACKOFF_MAX_SEC,default:60,min:10,max:3600"`

	APIRecoveryIntervalSec int `env:"API_RECOVERY_INTERVAL_SEC,default:300,min:300,max:3600"`
	APIRecoveryTimeoutSec  int `env:"API_RECOVERY_TIMEOUT_SEC,default:15,min:5,max:60"`
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
