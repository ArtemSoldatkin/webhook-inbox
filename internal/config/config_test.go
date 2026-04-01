package config

import (
	"os"
	"os/exec"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestLoadConfig(t *testing.T) {
	setValidConfigEnv(t)
	t.Setenv("API_DEFAULT_PAGE_SIZE", "25")
	t.Setenv("API_MIN_PAGE_SIZE", "5")
	t.Setenv("API_MAX_PAGE_SIZE", "50")
	t.Setenv("API_DELIVERY_RETRY_BACKOFF_BASE_SEC", "2")
	t.Setenv("API_DELIVERY_RETRY_BACKOFF_MAX_SEC", "20")

	cfg := LoadConfig()

	assert.Equal(t, "dev", cfg.Env)
	assert.Equal(t, "postgres", cfg.DBUser)
	assert.Equal(t, "secret", cfg.DBPassword)
	assert.Equal(t, "localhost", cfg.DBHost)
	assert.Equal(t, 5432, cfg.DBPort)
	assert.Equal(t, "webhook_inbox", cfg.DBName)
	assert.Equal(t, "http", cfg.APIProtocol)
	assert.Equal(t, "127.0.0.1", cfg.APIHost)
	assert.Equal(t, 8080, cfg.APIPort)
	assert.Equal(t, 25, cfg.APIDefaultPageSize)
	assert.Equal(t, 5, cfg.APIMinPageSize)
	assert.Equal(t, 50, cfg.APIMaxPageSize)
	assert.Equal(t, 2, cfg.APIDeliveryRetryBackoffBaseSec)
	assert.Equal(t, 20, cfg.APIDeliveryRetryBackoffMaxSec)
	assert.Equal(t, "http", cfg.UIProtocol)
	assert.Equal(t, "localhost", cfg.UIHost)
	assert.Equal(t, 3000, cfg.UIPort)
}

func TestLoadConfig_FatalScenarios(t *testing.T) {
	testCases := []struct {
		name           string
		extraEnv       map[string]string
		expectedOutput string
	}{
		{
			name: "min page greater than max page",
			extraEnv: map[string]string{
				"API_MIN_PAGE_SIZE": "20",
				"API_MAX_PAGE_SIZE": "10",
			},
			expectedOutput: "API_MIN_PAGE_SIZE cannot be greater than API_MAX_PAGE_SIZE",
		},
		{
			name: "default page outside bounds",
			extraEnv: map[string]string{
				"API_DEFAULT_PAGE_SIZE": "20",
				"API_MIN_PAGE_SIZE":     "30",
				"API_MAX_PAGE_SIZE":     "40",
			},
			expectedOutput: "API_DEFAULT_PAGE_SIZE must be between API_MIN_PAGE_SIZE and API_MAX_PAGE_SIZE",
		},
		{
			name: "retry backoff base greater than max",
			extraEnv: map[string]string{
				"API_DELIVERY_RETRY_BACKOFF_BASE_SEC": "30",
				"API_DELIVERY_RETRY_BACKOFF_MAX_SEC":  "20",
			},
			expectedOutput: "API_DELIVERY_RETRY_BACKOFF_BASE_SEC cannot be greater than API_DELIVERY_RETRY_BACKOFF_MAX_SEC",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cmd := exec.Command(os.Args[0], "-test.run=TestLoadConfigFatalHelper")
			cmd.Env = append(baseValidConfigEnv(), "CONFIG_FATAL_HELPER=1")
			for key, value := range tc.extraEnv {
				cmd.Env = append(cmd.Env, key+"="+value)
			}

			output, err := cmd.CombinedOutput()

			require.Error(t, err)
			exitErr, ok := err.(*exec.ExitError)
			require.True(t, ok)
			assert.NotEqual(t, 0, exitErr.ExitCode())
			assert.Contains(t, string(output), tc.expectedOutput)
		})
	}
}

func TestLoadConfigFatalHelper(t *testing.T) {
	if os.Getenv("CONFIG_FATAL_HELPER") != "1" {
		return
	}

	LoadConfig()
	t.Fatalf("expected LoadConfig to exit")
}

func setValidConfigEnv(t *testing.T) {
	t.Helper()
	for _, entry := range baseValidConfigEnv() {
		key, value, _ := splitEnv(entry)
		t.Setenv(key, value)
	}
}

func baseValidConfigEnv() []string {
	return []string{
		"ENV=dev",
		"POSTGRES_USER=postgres",
		"POSTGRES_PASSWORD=secret",
		"POSTGRES_HOST=localhost",
		"POSTGRES_PORT=5432",
		"POSTGRES_DB=webhook_inbox",
		"API_PROTOCOL=http",
		"API_HOST=127.0.0.1",
		"API_PORT=8080",
		"API_RATE_LIMIT_REQUESTS=100",
		"API_RATE_LIMIT_WINDOW_SEC=60",
		"API_THROTTLE_CONCURRENT_LIMIT=100",
		"API_REQUEST_SIZE_LIMIT_BYTES=10485760",
		"API_CORS_MAX_AGE_SEC=300",
		"API_CACHE_NUM_COUNTERS=10000000",
		"API_CACHE_MAX_COST=1073741824",
		"API_CACHE_BUFFER_ITEMS=64",
		"API_CACHE_DEFAULT_COST=1024",
		"API_CACHE_SOURCE_TTL_SEC=300",
		"API_CACHE_EVENT_TTL_SEC=900",
		"API_DEFAULT_PAGE_SIZE=20",
		"API_MIN_PAGE_SIZE=1",
		"API_MAX_PAGE_SIZE=100",
		"API_MAX_SEARCH_QUERY_LENGTH=512",
		"API_DELIVERY_INTERVAL_SEC=30",
		"API_DELIVERY_MAX_CONCURRENCY=10",
		"API_DELIVERY_TIMEOUT_SEC=15",
		"API_DELIVERY_REQUEST_TIMEOUT_SEC=15",
		"API_DELIVERY_MAX_RETRIES=3",
		"API_DELIVERY_RETRY_BACKOFF_BASE_SEC=1",
		"API_DELIVERY_RETRY_BACKOFF_MAX_SEC=60",
		"API_RECOVERY_INTERVAL_SEC=300",
		"API_RECOVERY_TIMEOUT_SEC=15",
		"UI_PROTOCOL=http",
		"UI_HOST=localhost",
		"UI_PORT=3000",
	}
}

func splitEnv(entry string) (string, string, bool) {
	for i := 0; i < len(entry); i++ {
		if entry[i] == '=' {
			return entry[:i], entry[i+1:], true
		}
	}
	return "", "", false
}
