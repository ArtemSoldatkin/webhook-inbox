package types

import (
	"encoding/base64"
	"testing"
	"time"

	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
)

func TestCursorToStringAndFromString_RoundTrip(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2026, 4, 1, 10, 20, 30, 123456789, time.FixedZone("UTC+2", 2*60*60))
	id := int64(42)
	cursor := NewCursor(&timestamp, &id)

	encoded := cursor.ToString()
	require.NotEmpty(t, encoded)

	var decoded Cursor
	require.NoError(t, decoded.FromString(encoded))

	require.NotNil(t, decoded.Timestamp)
	require.NotNil(t, decoded.ID)
	assert.Equal(t, timestamp.UTC(), decoded.Timestamp.UTC())
	assert.Equal(t, id, *decoded.ID)
	assert.Equal(t, encoded, decoded.String())
}

func TestCursorToString_ReturnsEmptyWhenMissingFields(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2026, 4, 1, 10, 20, 30, 0, time.UTC)
	id := int64(42)

	assert.Equal(t, "", (&Cursor{}).ToString())
	assert.Equal(t, "", (&Cursor{Timestamp: &timestamp}).ToString())
	assert.Equal(t, "", (&Cursor{ID: &id}).ToString())
}

func TestCursorFromString_EmptyInputResetsCursor(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2026, 4, 1, 10, 20, 30, 0, time.UTC)
	id := int64(99)
	cursor := Cursor{
		Timestamp: &timestamp,
		ID:        &id,
	}

	require.NoError(t, cursor.FromString(""))
	assert.Nil(t, cursor.Timestamp)
	assert.Nil(t, cursor.ID)
}

func TestCursorFromString_InvalidBase64(t *testing.T) {
	t.Parallel()

	var cursor Cursor
	err := cursor.FromString("%%%")

	require.Error(t, err)
}

func TestCursorFromString_InvalidFormat(t *testing.T) {
	t.Parallel()

	var cursor Cursor
	encoded := base64.URLEncoding.EncodeToString([]byte("2026-04-01T10:20:30Z"))

	err := cursor.FromString(encoded)

	require.Error(t, err)
	assert.Equal(t, "invalid cursor format", err.Error())
}

func TestCursorFromString_InvalidTimestamp(t *testing.T) {
	t.Parallel()

	var cursor Cursor
	encoded := base64.URLEncoding.EncodeToString([]byte("not-a-time|42"))

	err := cursor.FromString(encoded)

	require.Error(t, err)
	assert.Equal(t, "invalid cursor timestamp parameter - not-a-time", err.Error())
}

func TestCursorFromString_InvalidID(t *testing.T) {
	t.Parallel()

	var cursor Cursor
	encoded := base64.URLEncoding.EncodeToString([]byte("2026-04-01T10:20:30Z|abc"))

	err := cursor.FromString(encoded)

	require.Error(t, err)
	assert.Equal(t, "invalid cursor id parameter - abc", err.Error())
}

func TestCursorToDBParams(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2026, 4, 1, 10, 20, 30, 0, time.UTC)
	id := int64(123)
	cursor := Cursor{
		Timestamp: &timestamp,
		ID:        &id,
	}

	dbTimestamp, dbID := cursor.ToDBParams()

	assert.True(t, dbTimestamp.Valid)
	assert.Equal(t, timestamp, dbTimestamp.Time)
	assert.Equal(t, int64(123), dbID)
}

func TestCursorToDBParams_UsesDefaultsWhenNil(t *testing.T) {
	t.Parallel()

	dbTimestamp, dbID := (&Cursor{}).ToDBParams()

	assert.False(t, dbTimestamp.Valid)
	assert.Equal(t, int64(-1), dbID)
}

func TestNewCursor(t *testing.T) {
	t.Parallel()

	timestamp := time.Date(2026, 4, 1, 10, 20, 30, 0, time.UTC)
	id := int64(7)

	cursor := NewCursor(&timestamp, &id)

	require.NotNil(t, cursor.Timestamp)
	require.NotNil(t, cursor.ID)
	assert.Equal(t, timestamp, *cursor.Timestamp)
	assert.Equal(t, id, *cursor.ID)
}

func TestParseCursorTimestamp(t *testing.T) {
	t.Parallel()

	tsRFC3339Nano, err := parseCursorTimestamp("2026-04-01T10:20:30.123456789Z")
	require.NoError(t, err)
	require.NotNil(t, tsRFC3339Nano)
	assert.Equal(t, time.Date(2026, 4, 1, 10, 20, 30, 123456789, time.UTC), *tsRFC3339Nano)

	tsRFC3339, err := parseCursorTimestamp("2026-04-01T10:20:30Z")
	require.NoError(t, err)
	require.NotNil(t, tsRFC3339)
	assert.Equal(t, time.Date(2026, 4, 1, 10, 20, 30, 0, time.UTC), *tsRFC3339)

	tsEmpty, err := parseCursorTimestamp("")
	require.NoError(t, err)
	assert.Nil(t, tsEmpty)
}

func TestParseCursorTimestamp_Invalid(t *testing.T) {
	t.Parallel()

	ts, err := parseCursorTimestamp("bad-timestamp")

	assert.Nil(t, ts)
	require.Error(t, err)
	assert.Equal(t, "invalid cursor timestamp parameter - bad-timestamp", err.Error())
}

func TestParseCursorID(t *testing.T) {
	t.Parallel()

	id, err := parseCursorID("55")
	require.NoError(t, err)
	require.NotNil(t, id)
	assert.Equal(t, int64(55), *id)

	emptyID, err := parseCursorID("")
	require.NoError(t, err)
	assert.Nil(t, emptyID)
}

func TestParseCursorID_Invalid(t *testing.T) {
	t.Parallel()

	id, err := parseCursorID("bad-id")

	assert.Nil(t, id)
	require.Error(t, err)
	assert.Equal(t, "invalid cursor id parameter - bad-id", err.Error())
}
