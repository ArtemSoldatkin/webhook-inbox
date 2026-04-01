package mapperv1

import (
	"net/netip"
	"testing"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestToEventDTO(t *testing.T) {
	t.Parallel()

	receivedAt := time.Date(2026, 2, 3, 4, 5, 6, 0, time.UTC)
	remoteAddress := netip.MustParseAddr("192.0.2.10")

	event := db.Event{
		ID:              101,
		SourceID:        202,
		DedupHash:       pgtype.Text{String: "dedup-hash", Valid: true},
		Method:          "POST",
		IngressPath:     "/ingest/source-1",
		RemoteAddress:   &remoteAddress,
		QueryParams:     []byte(`{"foo":["bar","baz"]}`),
		RawHeaders:      []byte(`{"X-Test":["value"]}`),
		Body:            []byte(`{"ok":true}`),
		BodyContentType: "application/json",
		ReceivedAt:      pgtype.Timestamptz{Time: receivedAt, Valid: true},
	}

	dto := ToEventDTO(event)

	require.NotNil(t, dto.RemoteAddress)
	assert.Equal(t, int64(101), dto.ID)
	assert.Equal(t, int64(202), dto.SourceID)
	assert.Equal(t, "dedup-hash", dto.DedupHash)
	assert.Equal(t, "POST", dto.Method)
	assert.Equal(t, "/ingest/source-1", dto.IngressPath)
	assert.Equal(t, "192.0.2.10", *dto.RemoteAddress)
	assert.Equal(t, map[string][]string{"foo": {"bar", "baz"}}, dto.QueryParams)
	assert.Equal(t, map[string][]string{"X-Test": {"value"}}, dto.RawHeaders)
	assert.Equal(t, []byte(`{"ok":true}`), dto.Body)
	assert.Equal(t, "application/json", dto.BodyContentType)
	assert.Equal(t, receivedAt, dto.ReceivedAt)
}

func TestToEventDTO_InvalidJSONFallsBack(t *testing.T) {
	t.Parallel()

	dto := ToEventDTO(db.Event{
		QueryParams: []byte(`{`),
		RawHeaders:  []byte(`not-json`),
	})

	assert.Nil(t, dto.RemoteAddress)
	assert.Equal(t, map[string][]string{
		"__error": {"Webhook Inbox Error - Failed to parse"},
	}, dto.QueryParams)
	assert.Equal(t, map[string][]string{
		"__error": {"Webhook Inbox Error - Failed to parse"},
	}, dto.RawHeaders)
}

func TestToEventDTOs(t *testing.T) {
	t.Parallel()

	events := []db.Event{
		{ID: 1, QueryParams: []byte(`{"a":["1"]}`), RawHeaders: []byte(`{"H":["v"]}`)},
		{ID: 2, QueryParams: []byte(`{"b":["2"]}`), RawHeaders: []byte(`{"K":["w"]}`)},
	}

	dtos := ToEventDTOs(events)

	require.Len(t, dtos, 2)
	assert.Equal(t, int64(1), dtos[0].ID)
	assert.Equal(t, map[string][]string{"a": {"1"}}, dtos[0].QueryParams)
	assert.Equal(t, int64(2), dtos[1].ID)
	assert.Equal(t, map[string][]string{"b": {"2"}}, dtos[1].QueryParams)
}
