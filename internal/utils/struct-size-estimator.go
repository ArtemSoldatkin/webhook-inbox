package utils

import (
	"errors"
	"fmt"
	"net/netip"
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
)

const stringHeaderSize = 16
const sliceHeaderSize = 16

// EstimateStructSize estimates the size of a struct in bytes based on its fields.
// It supports basic types like int, string, byte slices, and some pgtype types.
// The function iterates through the fields of the struct and calculates the size based on the type of each field.
// If a field type is unsupported, it returns an error.
func EstimateStructSize(s interface{}) (int64, error) {
	if reflect.TypeOf(s).Kind() != reflect.Struct {
		return 0, errors.New("input must be a struct")
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
		return int64(field.Type.Size()), nil
	case reflect.TypeOf(""):
		return stringHeaderSize + int64(len(value.String())), nil
	case reflect.TypeOf([]byte{}):
		return sliceHeaderSize + int64(len(value.Bytes())), nil
	case reflect.TypeOf(pgtype.Text{}):
		// Only estimate the string inside, not the struct overhead
		return stringHeaderSize + int64(len(value.Interface().(pgtype.Text).String)), nil
	case reflect.TypeOf(pgtype.UUID{}):
		return int64(field.Type.Size()), nil // struct size, usually 16 bytes
	case reflect.TypeOf(&netip.Addr{}):
		return int64(field.Type.Size()), nil // pointer size, usually 8 bytes
	case reflect.TypeOf(pgtype.Timestamptz{}):
		return int64(field.Type.Size()), nil // struct size, usually 8 bytes
	default:
		return 0, fmt.Errorf(
			"unsupported field type %s for field %s",
			field.Type,
			field.Name,
		)
	}
}
