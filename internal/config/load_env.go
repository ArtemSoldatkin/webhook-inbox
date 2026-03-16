package config

import (
	"os"

	structparser "github.com/ArtemSoldatkin/webhook-inbox/internal/struct_parser"
)

// loadEnvs is a helper function that uses reflection to load environment variables into a struct based on struct field tags.
// It checks for required fields, allowed values, and default values as specified in the struct tags.
func loadEnvs[T any](config *T) error {
	return structparser.ParseStruct(
		config,
		"env",
		func(varName string) string {
			return os.Getenv(varName)
		},
		false,
	)
}
