package config

import (
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

type testEnvConfig struct {
	Required string `env:"TEST_REQUIRED,required"`
	Allowed  string `env:"TEST_ALLOWED,allowed:one|two"`
	Default  int    `env:"TEST_DEFAULT,default:7,min:1,max:10"`
}

func TestLoadEnvs(t *testing.T) {
	t.Setenv("TEST_REQUIRED", "configured")
	t.Setenv("TEST_ALLOWED", "two")
	t.Setenv("TEST_DEFAULT", "9")

	var cfg testEnvConfig
	err := loadEnvs(&cfg)

	require.NoError(t, err)
	assert.Equal(t, "configured", cfg.Required)
	assert.Equal(t, "two", cfg.Allowed)
	assert.Equal(t, 9, cfg.Default)
}

func TestLoadEnvs_AppliesDefaults(t *testing.T) {
	t.Setenv("TEST_REQUIRED", "configured")
	t.Setenv("TEST_ALLOWED", "one")

	var cfg testEnvConfig
	err := loadEnvs(&cfg)

	require.NoError(t, err)
	assert.Equal(t, 7, cfg.Default)
}

func TestLoadEnvs_ReturnsErrorForMissingRequired(t *testing.T) {
	t.Setenv("TEST_ALLOWED", "one")

	var cfg testEnvConfig
	err := loadEnvs(&cfg)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "TEST_REQUIRED")
}
