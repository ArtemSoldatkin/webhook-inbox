package utils

import (
	"errors"
	"fmt"
	"net/netip"
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
)

// EstimateStructSize estimates the size of a struct in bytes based on its fields.
// It supports basic types like int, string, byte slices, and some pgtype types.
// The function iterates through the fields of the struct and calculates the size based on the type of each field.
// If a field type is unsupported, it returns an error.
func EstimateStructSize(s interface{}) (int64, error) {
	if reflect.TypeOf(s).Kind() != reflect.Struct {
		return 0, errors.New("config must be a struct")
	}

	size := int64(0)
	for i := 0; i < reflect.TypeOf(s).NumField(); i++ {
		field := reflect.TypeOf(s).Field(i)
		value := reflect.ValueOf(s).Field(i)

		fieldSize, err := estimateFieldSize(field, value)
		if err != nil {
			return 0, err
		}

		size += fieldSize
	}

	return size, nil
}

// estimateFieldSize estimates the size of a single struct field based on its type and value.
func estimateFieldSize(
	field reflect.StructField,
	value reflect.Value,
) (int64, error) {
	switch field.Type {
	case reflect.TypeOf(0), reflect.TypeOf(int32(0)), reflect.TypeOf(int64(0)):
		return 8, nil
	case reflect.TypeOf(""):
		return int64(len(value.String())), nil
	case reflect.TypeOf([]byte{}):
		return int64(len(value.Bytes())), nil
	case reflect.TypeOf(pgtype.Text{}):
		return int64(len(value.Interface().(pgtype.Text).String)), nil
	case reflect.TypeOf(pgtype.UUID{}):
		return 16, nil // Size of UUID (16 bytes)
	case reflect.TypeOf(&netip.Addr{}):
		return 16, nil // Size of netip.Addr struct
	case reflect.TypeOf(pgtype.Timestamptz{}):
		return 8, nil // Size of timestamptz (8 bytes)
	default:
		return 0, fmt.Errorf(
			"unsupported field type %s for field %s",
			field.Type,
			field.Name,
		)
	}
}
