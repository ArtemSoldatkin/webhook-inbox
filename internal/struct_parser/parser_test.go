package structparser

import (
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

type parserConfig struct {
	IntValue    int          `env:"INT_VALUE,min:1,max:10"`
	Int64Value  int64        `env:"INT64_VALUE,required"`
	StringValue string       `env:"STRING_VALUE,allowed:one|two,max_length:5"`
	CursorValue types.Cursor `env:"CURSOR_VALUE"`
}

type missingTagConfig struct {
	Value string
}

type unsupportedConfig struct {
	Flag bool `env:"FLAG"`
}

func TestParseStruct(t *testing.T) {
	t.Parallel()

	cursorTime := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	cursorID := int64(42)
	cursor := types.NewCursor(&cursorTime, &cursorID)

	var cfg parserConfig
	err := ParseStruct(&cfg, "env", func(name string) string {
		switch name {
		case "INT_VALUE":
			return "5"
		case "INT64_VALUE":
			return "99"
		case "STRING_VALUE":
			return "two"
		case "CURSOR_VALUE":
			return cursor.ToString()
		default:
			return ""
		}
	}, false)

	require.NoError(t, err)
	assert.Equal(t, 5, cfg.IntValue)
	assert.Equal(t, int64(99), cfg.Int64Value)
	assert.Equal(t, "two", cfg.StringValue)
	require.NotNil(t, cfg.CursorValue.Timestamp)
	require.NotNil(t, cfg.CursorValue.ID)
	assert.Equal(t, cursorTime, *cfg.CursorValue.Timestamp)
	assert.Equal(t, cursorID, *cfg.CursorValue.ID)
}

func TestParseStruct_IgnoresMissingTagWhenAllowed(t *testing.T) {
	t.Parallel()

	cfg := missingTagConfig{}

	err := ParseStruct(&cfg, "env", func(string) string { return "" }, true)

	require.NoError(t, err)
}

func TestParseStruct_ReturnsErrorForMissingTag(t *testing.T) {
	t.Parallel()

	cfg := missingTagConfig{}

	err := ParseStruct(&cfg, "env", func(string) string { return "" }, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "missing env tag for field Value")
}

func TestParseStruct_ReturnsErrorForUnsupportedFieldType(t *testing.T) {
	t.Parallel()

	cfg := unsupportedConfig{}

	err := ParseStruct(&cfg, "env", func(string) string { return "true" }, false)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported field type bool for field Flag")
}

func TestParseStruct_ReturnsErrorWhenConfigIsNotStruct(t *testing.T) {
	t.Parallel()

	value := 10

	err := ParseStruct(&value, "env", func(string) string { return "" }, false)

	require.Error(t, err)
	assert.Equal(t, "config must be a struct", err.Error())
}

func TestParseTag(t *testing.T) {
	t.Parallel()

	name, options, err := parseTag("FIELD,required,default:10,allowed:ASC|DESC")

	require.NoError(t, err)
	assert.Equal(t, "FIELD", name)
	assert.Equal(t, "", options["required"])
	assert.Equal(t, "10", options["default"])
	assert.Equal(t, "ASC|DESC", options["allowed"])
}

func TestGetVarWithDefault(t *testing.T) {
	t.Parallel()

	value, err := getVarWithDefault("FIELD", map[string]string{"default": "fallback"}, func(string) string {
		return ""
	})

	require.NoError(t, err)
	assert.Equal(t, "fallback", value)
}

func TestGetVarWithDefault_ReturnsErrorForMissingRequired(t *testing.T) {
	t.Parallel()

	value, err := getVarWithDefault("FIELD", map[string]string{"required": ""}, func(string) string {
		return ""
	})

	assert.Equal(t, "", value)
	require.Error(t, err)
	assert.Equal(t, "variable FIELD is required but not set", err.Error())
}

func TestSetIntField(t *testing.T) {
	t.Parallel()

	var cfg struct {
		Value int
	}
	field := reflectValue(&cfg.Value)

	err := setIntField(field, "VALUE", "Value", map[string]string{"min": "1", "max": "10"}, strconv.IntSize, func(string) string {
		return "7"
	})

	require.NoError(t, err)
	assert.Equal(t, 7, cfg.Value)
}

func TestSetIntField_ReturnsErrorForInvalidValue(t *testing.T) {
	t.Parallel()

	var cfg struct {
		Value int
	}

	err := setIntField(reflectValue(&cfg.Value), "VALUE", "Value", nil, strconv.IntSize, func(string) string {
		return "bad"
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "variable VALUE has invalid value for field Value")
}

func TestSetStringField(t *testing.T) {
	t.Parallel()

	var cfg struct {
		Value string
	}

	err := setStringField(reflectValue(&cfg.Value), "VALUE", "Value", map[string]string{"allowed": "one|two", "max_length": "3"}, func(string) string {
		return "two"
	})

	require.NoError(t, err)
	assert.Equal(t, "two", cfg.Value)
}

func TestSetStringField_ReturnsErrorForDisallowedValue(t *testing.T) {
	t.Parallel()

	var cfg struct {
		Value string
	}

	err := setStringField(reflectValue(&cfg.Value), "VALUE", "Value", map[string]string{"allowed": "one|two"}, func(string) string {
		return "three"
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "variable VALUE has invalid value for field Value")
}

func TestIsIntValueWithinBoundary(t *testing.T) {
	t.Parallel()

	assert.True(t, isIntValueWithinBoundary(5, map[string]string{"min": "1", "max": "10"}, 64))
	assert.False(t, isIntValueWithinBoundary(0, map[string]string{"min": "1"}, 64))
	assert.False(t, isIntValueWithinBoundary(11, map[string]string{"max": "10"}, 64))
	assert.False(t, isIntValueWithinBoundary(5, map[string]string{"min": "bad"}, 64))
}

func TestIsStringValueAllowed(t *testing.T) {
	t.Parallel()

	assert.True(t, isStringValueAllowed("one", map[string]string{"allowed": "one|two"}))
	assert.False(t, isStringValueAllowed("three", map[string]string{"allowed": "one|two"}))
	assert.True(t, isStringValueAllowed("anything", nil))
}

func TestIsStringValueWithinBoundary(t *testing.T) {
	t.Parallel()

	assert.True(t, isStringValueWithinBoundary("test", map[string]string{"min_length": "2", "max_length": "5"}))
	assert.False(t, isStringValueWithinBoundary("a", map[string]string{"min_length": "2"}))
	assert.False(t, isStringValueWithinBoundary("toolong", map[string]string{"max_length": "5"}))
	assert.False(t, isStringValueWithinBoundary("test", map[string]string{"max_length": "bad"}))
}

func TestSetCursorField(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
	id := int64(33)
	cursor := types.NewCursor(&timestamp, &id)
	cursorString := cursor.ToString()

	var cfg struct {
		Cursor types.Cursor
	}

	err := setCursorField(reflectValue(&cfg.Cursor), "CURSOR", "Cursor", nil, func(string) string {
		return cursorString
	})

	require.NoError(t, err)
	require.NotNil(t, cfg.Cursor.Timestamp)
	require.NotNil(t, cfg.Cursor.ID)
	assert.Equal(t, timestamp, *cfg.Cursor.Timestamp)
	assert.Equal(t, id, *cfg.Cursor.ID)
}

func TestSetCursorField_ReturnsErrorForInvalidCursor(t *testing.T) {
	t.Parallel()

	var cfg struct {
		Cursor types.Cursor
	}

	err := setCursorField(reflectValue(&cfg.Cursor), "CURSOR", "Cursor", nil, func(string) string {
		return "bad-cursor"
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "error parsing cursor for field Cursor")
}

func reflectValue[T any](value *T) reflect.Value {
	return reflect.ValueOf(value).Elem()
}
