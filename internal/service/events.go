package service

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/netip"
	"time"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/api/types"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/sirupsen/logrus"
)

// ListEvents retrieves all events for a given source ID from the database.
func (svc *Service) ListEvents(
	ctx context.Context,
	sourceID int64,
	cursor types.Cursor,
	pageSize int,
	searchQuery string,
	sortDirection string,
) ([]db.Event, error) {
	cursorTS, cursorID := cursor.ToDBParams()
	return svc.queries.ListEventsBySource(ctx, db.ListEventsBySourceParams{
		SourceID:      sourceID,
		CursorTs:      cursorTS,
		CursorID:      cursorID,
		SearchQuery:   searchQuery,
		PageSize:      int32(pageSize),
		SortDirection: sortDirection,
	})
}

// GetEventByID retrieves a specific event by its ID from the database.
func (svc *Service) GetEventByID(ctx context.Context, eventID int64) (db.Event, error) {
	cacheKey := fmt.Sprintf("GetEventByID|%d", eventID)

	if cachedEvent, ok := svc.Cache.Get(cacheKey); ok {
		event, ok := cachedEvent.(db.Event)
		if ok {
			return event, nil
		}

		logrus.
			WithField("event_id", eventID).
			Warning("cache hit for GetEventByID but value has unexpected type, ignoring cache")
	}

	event, err := svc.queries.GetEventByID(ctx, eventID)
	if err != nil {
		return db.Event{}, err
	}

	cacheCost, err := utils.EstimateStructSize(event)
	if err != nil {
		logrus.
			WithError(err).
			WithField("event_id", event.ID).
			Warning("failed to estimate cache cost for event, using default cost")
		cacheCost = svc.Config.APICacheDefaultCost
	}

	svc.Cache.SetWithTTL(
		cacheKey,
		event,
		cacheCost,
		time.Duration(svc.Config.APICacheEventTTLSec)*time.Second,
	)

	return event, nil
}

// SourceIsNotFound represents an error when a source
// with the given public ID is not found in the database.
type SourceIsNotFound struct {
	Message string `json:"error"`
	Err     error  `json:"-"`
}

// Error returns the error message for SourceIsNotFound.
func (e *SourceIsNotFound) Error() string {
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work with SourceIsNotFound.
func (e *SourceIsNotFound) Unwrap() error { return e.Err }

// CreateEvent inserts a new event into the database and returns its ID.
func (svc *Service) CreateEvent(
	ctx context.Context,
	r *http.Request,
	publicID string,
) (int64, error) {
	source, err := svc.GetSourceByPublicID(r.Context(), publicID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, &SourceIsNotFound{
				Message: fmt.Sprintf("source with public_id '%s' not found", publicID),
				Err:     err,
			}
		}
		return 0, err
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return 0, fmt.Errorf("failed to split host and port from remote address: %w", err)
	}

	remoteAddress, err := netip.ParseAddr(host)
	if err != nil {
		return 0, fmt.Errorf("failed to parse remote address: %w", err)
	}

	headerBytes, err := json.Marshal(r.Header)
	if err != nil {
		return 0, err
	}

	queryParams, err := json.Marshal(r.URL.Query())
	if err != nil {
		return 0, err
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, err
	}
	r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

	method := r.Method
	ingressPath := r.URL.Path

	dedupHash, err := generateDedupHash(DedupPayload{
		Method:      method,
		IngressPath: ingressPath,
		QueryParams: queryParams,
		RawHeaders:  headerBytes,
		Body:        bodyBytes,
	})
	if err != nil {
		return 0, err
	}

	return svc.queries.CreateEvent(ctx,
		db.CreateEventParams{
			SourceID:        source.ID,
			DedupHash:       pgtype.Text{String: dedupHash, Valid: dedupHash != ""},
			Method:          method,
			IngressPath:     ingressPath,
			RemoteAddress:   &remoteAddress,
			QueryParams:     queryParams,
			RawHeaders:      headerBytes,
			Body:            bodyBytes,
			BodyContentType: r.Header.Get("Content-Type"),
		})
}

// DedupPayload represents the data used to generate a deduplication hash for an event.
type DedupPayload struct {
	Method      string `json:"method"`
	IngressPath string `json:"ingress_path"`
	QueryParams []byte `json:"query_params"`
	RawHeaders  []byte `json:"headers"`
	Body        []byte `json:"body"`
}

// generateDedupHash generates a SHA-256 hash of the input data for deduplication purposes.
func generateDedupHash(dedupPayload DedupPayload) (string, error) {
	dedupData, err := json.Marshal(dedupPayload)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(dedupData)
	return hex.EncodeToString(hash[:]), nil
}
