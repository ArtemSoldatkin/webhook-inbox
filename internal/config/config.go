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

	APIRateLimitRequests       int   `env:"API_RATE_LIMIT_REQUESTS,default:100,min:1,max:1000"`
	APIRateLimitWindowSec      int   `env:"API_RATE_LIMIT_WINDOW_SEC,default:60,min:1,max:3600"`
	APIThrottleConcurrentLimit int   `env:"API_THROTTLE_CONCURRENT_LIMIT,default:100,min:1,max:1000"`
	APIRequestSizeLimitBytes   int64 `env:"API_REQUEST_SIZE_LIMIT_BYTES,default:10485760,min:1024,max:104857600"`

	APICORSMaxAgeSec int `env:"API_CORS_MAX_AGE_SEC,default:300,min:60,max:3600"`

	APICacheNumCounters int64 `env:"API_CACHE_NUM_COUNTERS,default:10000000,min:1000,max:100000000"`
	APICacheMaxCost     int64 `env:"API_CACHE_MAX_COST,default:1073741824,min:1048576,max:10737418240"`
	APICacheBufferItems int64 `env:"API_CACHE_BUFFER_ITEMS,default:64,min:1,max:1024"`

	APICacheDefaultCost int64 `env:"API_CACHE_DEFAULT_COST,default:1024,min:1,max:1048576"`

	APICacheSourceTTLSec int `env:"API_CACHE_SOURCE_TTL_SEC,default:300,min:60,max:3600"`
	APICacheEventTTLSec  int `env:"API_CACHE_EVENT_TTL_SEC,default:900,min:60,max:3600"`

	APIDefaultPageSize int `env:"API_DEFAULT_PAGE_SIZE,default:20,min:1,max:100"`
	APIMinPageSize     int `env:"API_MIN_PAGE_SIZE,default:1,min:1,max:100"`
	APIMaxPageSize     int `env:"API_MAX_PAGE_SIZE,default:100,min:1,max:100"`

	APIMaxSearchQueryLength int `env:"API_MAX_SEARCH_QUERY_LENGTH,default:512,min:1,max:1024"`

	APIDeliveryIntervalSec       int `env:"API_DELIVERY_INTERVAL_SEC,default:30,min:10,max:60"`
	APIDeliveryMaxConcurrency    int `env:"API_DELIVERY_MAX_CONCURRENCY,default:10,min:1,max:100"`
	APIDeliveryTimeoutSec        int `env:"API_DELIVERY_TIMEOUT_SEC,default:15,min:5,max:60"`
	APIDeliveryRequestTimeoutSec int `env:"API_DELIVERY_REQUEST_TIMEOUT_SEC,default:15,min:5,max:60"`

	APIDeliveryMaxRetries          int `env:"API_DELIVERY_MAX_RETRIES,default:3,min:0,max:10"`
	APIDeliveryRetryBackoffBaseSec int `env:"API_DELIVERY_RETRY_BACKOFF_BASE_SEC,default:1,min:1,max:60"`
	APIDeliveryRetryBackoffMaxSec  int `env:"API_DELIVERY_RETRY_BACKOFF_MAX_SEC,default:60,min:10,max:3600"`

	APIRecoveryIntervalSec int `env:"API_RECOVERY_INTERVAL_SEC,default:300,min:300,max:3600"`
	APIRecoveryTimeoutSec  int `env:"API_RECOVERY_TIMEOUT_SEC,default:15,min:5,max:60"`

	UIProtocol string `env:"UI_PROTOCOL,required,allowed:http|https"`
	UIHost     string `env:"UI_HOST,required"`
	UIPort     int    `env:"UI_PORT,required"`
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() Config {
	var config Config
	err := loadEnvs(&config)
	if err != nil {
		logrus.
			WithError(err).
			Fatal("Error loading configuration")
	}

	if config.APIMinPageSize > config.APIMaxPageSize {
		logrus.Fatal("API_MIN_PAGE_SIZE cannot be greater than API_MAX_PAGE_SIZE")
	}

	if config.APIDefaultPageSize < config.APIMinPageSize ||
		config.APIDefaultPageSize > config.APIMaxPageSize {
		logrus.Fatal("API_DEFAULT_PAGE_SIZE must be between API_MIN_PAGE_SIZE and API_MAX_PAGE_SIZE")
	}

	if config.APIDeliveryRetryBackoffBaseSec > config.APIDeliveryRetryBackoffMaxSec {
		logrus.Fatal("API_DELIVERY_RETRY_BACKOFF_BASE_SEC cannot be greater than API_DELIVERY_RETRY_BACKOFF_MAX_SEC")
	}

	return config
}
