package utils

import (
	"net/netip"
	"reflect"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5/pgtype"
)

type supportedStruct struct {
	IntValue       int
	Int32Value     int32
	Int64Value     int64
	StringValue    string
	BytesValue     []byte
	TextValue      pgtype.Text
	UUIDValue      pgtype.UUID
	RemoteAddress  *netip.Addr
	TimestampValue pgtype.Timestamptz
}

type unsupportedStruct struct {
	Flag bool
}

func TestEstimateStructSize(t *testing.T) {
	t.Parallel()

	remoteAddress := netip.MustParseAddr("192.0.2.10")
	value := supportedStruct{
		IntValue:       1,
		Int32Value:     2,
		Int64Value:     3,
		StringValue:    "hello",
		BytesValue:     []byte("world"),
		TextValue:      pgtype.Text{String: "text", Valid: true},
		UUIDValue:      pgtype.UUID{},
		RemoteAddress:  &remoteAddress,
		TimestampValue: pgtype.Timestamptz{},
	}

	size, err := EstimateStructSize(value)

	require.NoError(t, err)

	expected := int64(reflect.TypeOf(value.IntValue).Size()) +
		int64(reflect.TypeOf(value.Int32Value).Size()) +
		int64(reflect.TypeOf(value.Int64Value).Size()) +
		int64(stringHeaderSize+len(value.StringValue)) +
		int64(sliceHeaderSize+len(value.BytesValue)) +
		int64(stringHeaderSize+len(value.TextValue.String)) +
		int64(reflect.TypeOf(value.UUIDValue).Size()) +
		int64(reflect.TypeOf(value.RemoteAddress).Size()) +
		int64(reflect.TypeOf(value.TimestampValue).Size())
	assert.Equal(t, expected, size)
}

func TestEstimateStructSize_NonStruct(t *testing.T) {
	t.Parallel()

	size, err := EstimateStructSize("not-a-struct")

	assert.Equal(t, int64(0), size)
	require.Error(t, err)
	assert.Equal(t, "input must be a struct", err.Error())
}

func TestEstimateStructSize_UnsupportedField(t *testing.T) {
	t.Parallel()

	size, err := EstimateStructSize(unsupportedStruct{Flag: true})

	assert.Equal(t, int64(0), size)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported field type bool for field Flag")
}

func TestEstimateFieldSize(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeOf(supportedStruct{})
	val := reflect.ValueOf(supportedStruct{
		StringValue: "hello",
	})
	field := typ.Field(3)
	fieldValue := val.Field(3)

	size, err := estimateFieldSize(field, fieldValue)

	require.NoError(t, err)
	assert.Equal(t, int64(stringHeaderSize+5), size)
}

func TestEstimateFieldSize_UnsupportedField(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeOf(unsupportedStruct{})
	val := reflect.ValueOf(unsupportedStruct{Flag: true})
	field := typ.Field(0)
	fieldValue := val.Field(0)

	size, err := estimateFieldSize(field, fieldValue)

	assert.Equal(t, int64(0), size)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported field type bool for field Flag")
}
