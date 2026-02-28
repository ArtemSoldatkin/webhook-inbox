package config

import (
	"errors"
	"fmt"
	"strconv"

	"os"
	"reflect"
	"strings"
)

// loadEnvs is a helper function that uses reflection to load environment variables into a struct based on struct field tags.
// It checks for required fields, allowed values, and default values as specified in the struct tags.
func loadEnvs[T any](config *T) error {
	if reflect.TypeOf(*config).Kind() != reflect.Struct {
		return errors.New("config must be a struct")
	}
	for i := 0; i < reflect.TypeOf(*config).NumField(); i++ {
		field := reflect.TypeOf(*config).Field(i)
		envTag := field.Tag.Get("env")
		if envTag == "" {
			return fmt.Errorf(
				"missing env tag for field %s",
				field.Name,
			)
		}
		envName, envOptions, err := parseTag(envTag)
		if err != nil {
			return fmt.Errorf(
				"error parsing env tag for field %s: %w",
				field.Name,
				err,
			)
		}
		configValue := reflect.ValueOf(config).Elem().Field(i)
		switch configValue.Kind() {
			case reflect.Int:
				if err := setIntField(
					configValue,
					envName,
					field.Name,
					envOptions,
				); err != nil {
					return err
				}
			case reflect.String:
				if err := setStringField(
					configValue,
					envName,
					field.Name,
					envOptions,
				); err != nil {
					return err
				}
			default:
				return fmt.Errorf(
					"unsupported field type %s for field %s",
					configValue.Kind(),
					reflect.TypeOf(*config).Field(i).Name,
				)
		}
	}
	return nil
}

// parseTag parses an `env` struct field tag value in the format "ENV_VAR,option1,option2:optionValue"
// and returns the environment variable name and a map of options.
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

// getEnvWithDefault retrieves the value of an environment variable,
// returning a default value if the variable is not set and a default is specified.
// It also checks if the variable is required and returns an error if it is not set.
func getEnvWithDefault(envName string, envOptions map[string]string) (string, error) {
	envValue := os.Getenv(envName)
	if envValue == "" {
		if defaultValue, hasDefault := envOptions["default"]; hasDefault {
			return defaultValue, nil
		}
		if _, hasRequired := envOptions["required"]; hasRequired {
			return "", fmt.Errorf(
				"environment variable %s is required but not set",
				envName,
			)
		}
	}
	return envValue, nil
}

// setIntField is a helper function that sets an integer field in the config struct based on the environment variable value.
// It retrieves the environment variable, checks for errors, converts it to an integer, and sets the field value.
func setIntField(configValue reflect.Value, envName string, fieldName string, envOptions map[string]string) error {
	envValue, err := getEnvWithDefault(envName, envOptions)
	if err != nil {
		return fmt.Errorf(
			"error getting environment variable %s for field %s: %w",
			envName,
			fieldName,
			err,
		)
	}
	valueInt, err := strconv.Atoi(envValue)
	if err != nil {
		return fmt.Errorf(
			"environment variable %s has invalid value for field %s: %w",
			envName,
			fieldName,
			err,
		)
	}
	if !isIntValueWithinBoundary(valueInt, envOptions) {
		return fmt.Errorf(
			"environment variable %s has value out of boundary for field %s: value '%d' is not within specified boundaries",
			envName,
			fieldName,
			valueInt,
		)
	}
	configValue.SetInt(int64(valueInt))
	return nil
}

// setStringField is a helper function that sets a string field in the config struct based on the environment variable value.
// It retrieves the environment variable, checks for errors, and sets the field value.
func setStringField(configValue reflect.Value, envName string, fieldName string, envOptions map[string]string) error {
	envValue, err := getEnvWithDefault(envName, envOptions)
	if err != nil {
		return fmt.Errorf(
			"error getting environment variable %s for field %s: %w",
			envName,
			fieldName,
			err,
		)
	}
	if !isStringValueAllowed(envValue, envOptions) {
		return fmt.Errorf(
			"environment variable %s has invalid value for field %s: value '%s' is not allowed",
			envName,
			fieldName,
			envValue,
		)
	}
	configValue.SetString(envValue)
	return nil
}

// isIntValueWithinBoundary checks if an integer value is within the specified minimum and maximum boundaries defined in the environment variable options.
// It returns true if the value is within the boundaries or if no boundaries are specified, and false otherwise.
func isIntValueWithinBoundary(value int, envOptions map[string]string) bool {
	minValueStr, hasMin := envOptions["min"]
	maxValueStr, hasMax := envOptions["max"]
	if hasMin && hasMax {
		minValue, err := strconv.Atoi(minValueStr)
		if err != nil {
			return false
		}
		maxValue, err := strconv.Atoi(maxValueStr)
		if err != nil {
			return false
		}
		return value >= minValue && value <= maxValue
	}
	if hasMin {
		minValue, err := strconv.Atoi(minValueStr)
		if err != nil {
			return false
		}
		return value >= minValue
	}
	if hasMax {
		maxValue, err := strconv.Atoi(maxValueStr)
		if err != nil {
			return false
		}
		return value <= maxValue
	}
	return true
}

// isStringValueAllowed checks if a string value is allowed based on the "allowed" option in the environment variable options.
// It returns true if the value is allowed or if there are no allowed values specified, and false otherwise.
func isStringValueAllowed(value string, envOptions map[string]string) bool {
	allowedValues, hasAllowed := envOptions["allowed"];
	if !hasAllowed {
		return true
	}
	allowedList := strings.Split(allowedValues, "|")
	for _, allowed := range allowedList {
		if value == allowed {
			return true
		}
	}
	return false
}
