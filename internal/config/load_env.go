package config

import (
	"errors"
	"fmt"
	"strconv"

	"os"
	"reflect"
	"strings"
)

// TODO Verify allowed values
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
		envName, _, err := parseTag(envTag)
		if err != nil {
			return fmt.Errorf("error parsing env tag for field %s: %w", field.Name, err)
		}
		configValue := reflect.ValueOf(config).Elem().Field(i)
		switch configValue.Kind() {
			case reflect.Int:
				if err := setIntField(configValue, envName, field.Name, nil); err != nil {
					return err
				}
			case reflect.String:
				if err := setStringField(configValue, envName, field.Name, nil); err != nil {
					return err
				}
			default:
				return fmt.Errorf("unsupported field type %s for field %s", configValue.Kind(), reflect.TypeOf(*config).Field(i).Name)
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

// getEnvWithDefault retrieves the value of an environment variable, returning a default value if the variable is not set and a default is specified. It also checks if the variable is required and returns an error if it is not set.
func getEnvWithDefault(envName string, envOptions map[string]string) (string, error) {
	envValue := os.Getenv(envName)
	if envValue == "" {
		if defaultValue, hasDefault := envOptions["default"]; hasDefault {
			return defaultValue, nil
		}
		if _, hasRequired := envOptions["required"]; hasRequired {
			return "", fmt.Errorf("environment variable %s is required but not set", envName)
		}
	}
	return envValue, nil
}

// setIntField is a helper function that sets an integer field in the config struct based on the environment variable value. It retrieves the environment variable, checks for errors, converts it to an integer, and sets the field value.
func setIntField(configValue reflect.Value, envName string, fieldName string, envOptions map[string]string) error {
	envValue, err := getEnvWithDefault(envName, envOptions)
	if err != nil {
		return fmt.Errorf("error getting environment variable %s for field %s: %w", envName, fieldName, err)
	}
	valueInt, err := strconv.Atoi(envValue)
	if err != nil {
		return fmt.Errorf("environment variable %s has invalid value for field %s: %w", envName, fieldName, err)
	}
	configValue.SetInt(int64(valueInt))
	return nil
}

// setStringField is a helper function that sets a string field in the config struct based on the environment variable value. It retrieves the environment variable, checks for errors, and sets the field value.
func setStringField(configValue reflect.Value, envName string, fieldName string, envOptions map[string]string) error {
	envValue, err := getEnvWithDefault(envName, envOptions)
	if err != nil {
		return fmt.Errorf("error getting environment variable %s for field %s: %w", envName, fieldName, err)
	}
	configValue.SetString(envValue)
	return nil
}
