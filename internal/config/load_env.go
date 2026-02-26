package config

import (
	"errors"
	"fmt"

	"os"
	"reflect"
	"strings"
)

// loadEnvs is a helper function that uses reflection to load environment variables into a struct based on struct field tags. It checks for required fields, allowed values, and default values as specified in the struct tags.
func loadEnvs[T any](config *T) error {
	if reflect.TypeOf(*config).Kind() != reflect.Struct {
		return errors.New("config must be a struct")
	}
	for i := 0; i < reflect.TypeOf(*config).NumField(); i++ {
		field := reflect.TypeOf(*config).Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			return fmt.Errorf("missing env tag for field %s", field.Name)
		}
		envName, envOptions, err := parseTag(envTag)
		if err != nil {
			return fmt.Errorf("error parsing env tag for field %s: %w", field.Name, err)
		}
		envValue := os.Getenv(envName)
		if envValue == "" {
			if _, required := envOptions["required"]; required {
				return fmt.Errorf("environment variable %s is required but not set", envName)
			}
			if defaultValue, hasDefault := envOptions["default"]; hasDefault {
				envValue = defaultValue
			}
		} else {
			if allowedValues, hasAllowed := envOptions["allowed"]; hasAllowed {
				allowedList := strings.Split(allowedValues, "|")
				isValid := false
				for _, allowed := range allowedList {
					if envValue == allowed {
						isValid = true
						break
					}
				}
				if !isValid {
					return fmt.Errorf("invalid value for environment variable %s: %s (allowed: %s)", envName, envValue, allowedValues)
				}
			} else {
				configValue := reflect.ValueOf(config).Elem().Field(i)
				if configValue.Kind() == reflect.Int {
					var intValue int
					_, err := fmt.Sscanf(envValue, "%d", &intValue)
					if err != nil {
						return fmt.Errorf("invalid integer value for environment variable %s: %s", envName, envValue)
					}
					configValue.SetInt(int64(intValue))
				} else if configValue.Kind() == reflect.String {
					configValue.SetString(envValue)
				} else {
					return fmt.Errorf("unsupported field type for environment variable %s: %s", envName, configValue.Kind())
				}
			}
		}
	}
	return nil
}

// parseTag parses a struct field tag in the format "env:ENV_VAR,option1,option2:optionValue" and returns the environment variable name and a map of options.
func parseTag(tag string) (string, map[string]string, error) {
	parts := strings.Split(tag, ",")
	if len(parts) < 1 {
		return "", nil, fmt.Errorf("invalid tag format: %s", tag)
	}
	envVar := parts[0]
	optionMap := make(map[string]string)
	for _, option := range parts[1:] {
		if strings.Contains(option, ":") {
			kv := strings.SplitN(option, ":", 2)
			optionMap[kv[0]] = kv[1]
		} else {
			optionMap[option] = ""
		}
	}
	return envVar, optionMap, nil
}