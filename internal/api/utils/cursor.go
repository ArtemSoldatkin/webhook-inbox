package api

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

// Cursor represents a pagination cursor with a timestamp and an ID,
// used for efficient pagination in API responses.
type Cursor struct {
	Timestamp *time.Time
	ID        *int64
}

// ToString encodes the Cursor into a base64 string for use in API responses.
func (c *Cursor) ToString() string {
	if c.Timestamp == nil || c.ID == nil {
		return ""
	}
	raw := fmt.Sprintf("%s|%d", c.Timestamp.UTC().Format(time.RFC3339Nano), *c.ID)
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

// FromString decodes a base64 string into a Cursor, parsing the timestamp and ID components.
func (c *Cursor) FromString(s string) error {
	if s == "" {
		c.Timestamp = nil
		c.ID = nil
		return nil
	}
	decoded, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return err
	}
	parts := strings.SplitN(string(decoded), "|", 2)
	if len(parts) != 2 {
		return errors.New("invalid cursor format")
	}
	ts, err := parseCursorTimestamp(parts[0])
	if err != nil {
		return err
	}
	id, err := parseCursorID(parts[1])
	if err != nil {
		return err
	}
	c.Timestamp = ts
	c.ID = id
	return nil
}

// ToDBParams converts the Cursor into database parameters (pgtype.Timestamptz and int64) for use in SQL queries.
func (c *Cursor) ToDBParams() (pgtype.Timestamptz, int64) {
	var ts pgtype.Timestamptz
	if c.Timestamp != nil {
		ts = pgtype.Timestamptz{Time: *c.Timestamp, Valid: true}
	} else {
		ts = pgtype.Timestamptz{Valid: false}
	}
	var id int64
	if c.ID != nil {
		id = *c.ID
	} else {
		id = -1 // Default to -1 if ID is nil, which should be treated as the maximum value for pagination
	}
	return ts, id
}

// NewCursor creates a new Cursor instance from the given timestamp and ID pointers.
func NewCursor(ts *time.Time, id *int64) Cursor {
	return Cursor{Timestamp: ts, ID: id}
}

// parseCursorTimestamp parses a string into a time.Time pointer,
// returning nil if the string is empty.
func parseCursorTimestamp(s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	parsedTime, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		parsedTime, err = time.Parse(time.RFC3339, s)
	}
	if err != nil {
		return nil, fmt.Errorf(
			"invalid cursor timestamp parameter - %s",
			s,
		)
	}
	return &parsedTime, nil
}

// parseCursorID parses a string into an int64 pointer, returning nil if the string is empty.
func parseCursorID(s string) (*int64, error) {
	if s == "" {
		return nil, nil
	}
	id, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return nil, fmt.Errorf(
			"invalid cursor id parameter - %s",
			s,
		)
	}
	return &id, nil
}
