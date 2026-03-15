package structparser

import (
	"errors"
	"fmt"
	"strconv"

	"reflect"
	"strings"
)

// ParseStruct is a generic function that populates a struct with values from variables based on struct tags
// It takes a pointer to the struct, the tag name to look for, and a function to retrieve variable values.
func ParseStruct[T any](config *T, tagName string, getVar func(string) string) error {
	if reflect.TypeOf(*config).Kind() != reflect.Struct {
		return errors.New("config must be a struct")
	}

	for i := 0; i < reflect.TypeOf(*config).NumField(); i++ {
		field := reflect.TypeOf(*config).Field(i)
		tag := field.Tag.Get(tagName)
		if tag == "" {
			return fmt.Errorf(
				"missing %s tag for field %s",
				tagName,
				field.Name,
			)
		}

		varName, varOptions, err := parseTag(tag)
		if err != nil {
			return fmt.Errorf(
				"error parsing %s tag for field %s: %w",
				tagName,
				field.Name,
				err,
			)
		}

		configValue := reflect.ValueOf(config).Elem().Field(i)
		switch configValue.Kind() {
		case reflect.Int:
			if err := setIntField(
				configValue,
				varName,
				field.Name,
				varOptions,
				strconv.IntSize,
				getVar,
			); err != nil {
				return err
			}
		case reflect.Int64:
			if err := setIntField(
				configValue,
				varName,
				field.Name,
				varOptions,
				64,
				getVar,
			); err != nil {
				return err
			}
		case reflect.String:
			if err := setStringField(
				configValue,
				varName,
				field.Name,
				varOptions,
				getVar,
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

// parseTag is a helper function that parses a struct tag into the variable name and options.
// It splits the tag by commas, extracts the variable name and options, and returns them.
func parseTag(tag string) (string, map[string]string, error) {
	parts := strings.Split(tag, ",")
	if len(parts) < 1 {
		return "", nil, fmt.Errorf("invalid tag format: %s", tag)
	}

	varName := parts[0]
	optionMap := make(map[string]string)
	for _, option := range parts[1:] {
		if strings.Contains(option, ":") {
			kv := strings.SplitN(option, ":", 2)
			optionMap[kv[0]] = kv[1]
		} else {
			optionMap[option] = ""
		}
	}

	return varName, optionMap, nil
}

// getVarWithDefault is a helper function that retrieves the variable value using the provided getVar function, and applies default and required logic based on the options.
// It checks if the variable is set, applies default value if not set, and returns an error if the variable is required but not set.
func getVarWithDefault(varName string, varOptions map[string]string, getVar func(string) string) (string, error) {
	varValue := getVar(varName)
	if varValue == "" {
		if defaultValue, hasDefault := varOptions["default"]; hasDefault {
			return defaultValue, nil
		}

		if _, hasRequired := varOptions["required"]; hasRequired {
			return "", fmt.Errorf(
				"variable %s is required but not set",
				varName,
			)
		}
	}

	return varValue, nil
}

// setIntField is a helper function that sets an integer field in the config struct based on the variable value.
// It retrieves the variable, checks for errors, validates boundaries, and sets the field value.
func setIntField(
	configValue reflect.Value,
	varName string,
	fieldName string,
	varOptions map[string]string,
	bitSize int,
	getVar func(string) string,
) error {
	varValue, err := getVarWithDefault(varName, varOptions, getVar)
	if err != nil {
		return fmt.Errorf(
			"error getting variable %s for field %s: %w",
			varName,
			fieldName,
			err,
		)
	}

	valueInt, err := strconv.ParseInt(varValue, 10, bitSize)
	if err != nil {
		return fmt.Errorf(
			"variable %s has invalid value for field %s: %w",
			varName,
			fieldName,
			err,
		)
	}

	if !isIntValueWithinBoundary(valueInt, varOptions, bitSize) {
		return fmt.Errorf(
			"variable %s has value out of boundary for field %s: value '%d' is not within specified boundaries",
			varName,
			fieldName,
			valueInt,
		)
	}

	configValue.SetInt(valueInt)
	return nil
}

// setStringField is a helper function that sets a string field in the config struct based on the variable value.
// It retrieves the variable, checks for errors, validates allowed values, and sets the field value.
func setStringField(
	configValue reflect.Value,
	varName string,
	fieldName string,
	varOptions map[string]string,
	getVar func(string) string,
) error {
	varValue, err := getVarWithDefault(varName, varOptions, getVar)
	if err != nil {
		return fmt.Errorf(
			"error getting variable %s for field %s: %w",
			varName,
			fieldName,
			err,
		)
	}

	if !isStringValueWithinBoundary(varValue, varOptions) {
		return fmt.Errorf(
			"variable %s has value length out of boundary for field %s: length of value '%s' is not within specified boundaries",
			varName,
			fieldName,
			varValue,
		)
	}

	if !isStringValueAllowed(varValue, varOptions) {
		return fmt.Errorf(
			"variable %s has invalid value for field %s: value '%s' is not allowed",
			varName,
			fieldName,
			varValue,
		)
	}

	configValue.SetString(varValue)
	return nil
}

// isIntValueWithinBoundary is a helper function that checks if an integer value is within the specified boundaries defined in the options.
// It checks for "min" and "max" options, parses them, and compares the value against the boundaries.
func isIntValueWithinBoundary(
	value int64,
	varOptions map[string]string,
	bitSize int,
) bool {
	minValueStr, hasMin := varOptions["min"]
	maxValueStr, hasMax := varOptions["max"]

	if hasMin && hasMax {
		minValue, err := strconv.ParseInt(minValueStr, 10, bitSize)
		if err != nil {
			return false
		}

		maxValue, err := strconv.ParseInt(maxValueStr, 10, bitSize)
		if err != nil {
			return false
		}

		return value >= minValue && value <= maxValue
	}

	if hasMin {
		minValue, err := strconv.ParseInt(minValueStr, 10, bitSize)
		if err != nil {
			return false
		}

		return value >= minValue
	}

	if hasMax {
		maxValue, err := strconv.ParseInt(maxValueStr, 10, bitSize)
		if err != nil {
			return false
		}

		return value <= maxValue
	}

	return true
}

// isStringValueAllowed is a helper function that checks if a string value is allowed based on the "allowed" option in the variable options.
// It splits the allowed values by "|" and checks if the value is in the list of allowed values.
func isStringValueAllowed(value string, varOptions map[string]string) bool {
	allowedValues, hasAllowed := varOptions["allowed"]
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

// isStringValueWithinBoundary is a helper function that checks if a string value's length is within the specified boundaries defined in the options.
// It checks for "min_length" and "max_length" options, parses them, and compares the length of the value against the boundaries.
func isStringValueWithinBoundary(
	value string,
	varOptions map[string]string,
) bool {
	minLengthValueStr, hasMin := varOptions["min_length"]
	maxLengthValueStr, hasMax := varOptions["max_length"]

	if hasMin && hasMax {
		minLengthValue, err := strconv.Atoi(minLengthValueStr)
		if err != nil {
			return false
		}

		maxLengthValue, err := strconv.Atoi(maxLengthValueStr)
		if err != nil {
			return false
		}

		return len(value) >= minLengthValue && len(value) <= maxLengthValue
	}

	if hasMin {
		minLengthValue, err := strconv.Atoi(minLengthValueStr)
		if err != nil {
			return false
		}

		return len(value) >= minLengthValue
	}

	if hasMax {
		maxLengthValue, err := strconv.Atoi(maxLengthValueStr)
		if err != nil {
			return false
		}

		return len(value) <= maxLengthValue
	}

	return true
}
