package routev1

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	dtov1 "github.com/ArtemSoldatkin/webhook-inbox/internal/api/dto/v1"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/service"
	"github.com/dgraph-io/ristretto"
	"github.com/go-chi/chi/v5"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func TestV1Router_MountsIngestAndSourcesRoutes(t *testing.T) {
	t.Parallel()

	router := V1Router(newTestService(t, newTestDB(), newTestConfig()))

	ingestReq := httptest.NewRequest("POST", "/ingest/not-a-uuid", nil)
	ingestRes := httptest.NewRecorder()
	router.ServeHTTP(ingestRes, ingestReq)

	assert.Equal(t, http.StatusBadRequest, ingestRes.Code)

	sourcesReq := httptest.NewRequest("GET", "/sources?limit=101", nil)
	sourcesRes := httptest.NewRecorder()
	router.ServeHTTP(sourcesRes, sourcesReq)

	assert.Equal(t, http.StatusBadRequest, sourcesRes.Code)
}

type paginatedSourcesResponse struct {
	Data       []dtov1.SourceDTO `json:"data"`
	NextCursor *string           `json:"next_cursor"`
	HasNext    bool              `json:"has_next"`
}

type paginatedEventsResponse struct {
	Data       []dtov1.EventDTO `json:"data"`
	NextCursor *string          `json:"next_cursor"`
	HasNext    bool             `json:"has_next"`
}

type paginatedDeliveryAttemptsResponse struct {
	Data       []dtov1.DeliveryAttemptDTO `json:"data"`
	NextCursor *string                    `json:"next_cursor"`
	HasNext    bool                       `json:"has_next"`
}

type testDB struct {
	queryHandlers    map[string]func(args ...any) (pgx.Rows, error)
	queryRowHandlers map[string]func(args ...any) pgx.Row
	execHandlers     map[string]func(args ...any) (pgconn.CommandTag, error)
}

func newTestDB() *testDB {
	return &testDB{
		queryHandlers:    make(map[string]func(args ...any) (pgx.Rows, error)),
		queryRowHandlers: make(map[string]func(args ...any) pgx.Row),
		execHandlers:     make(map[string]func(args ...any) (pgconn.CommandTag, error)),
	}
}

func (dbtx *testDB) Query(_ context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	handler, ok := dbtx.queryHandlers[querySignature(query)]
	if !ok {
		return nil, errors.New("unexpected query: " + querySignature(query))
	}
	return handler(args...)
}

func (dbtx *testDB) QueryRow(_ context.Context, query string, args ...interface{}) pgx.Row {
	handler, ok := dbtx.queryRowHandlers[querySignature(query)]
	if !ok {
		return &testRow{err: errors.New("unexpected query row: " + querySignature(query))}
	}
	return handler(args...)
}

func (dbtx *testDB) Exec(_ context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	handler, ok := dbtx.execHandlers[querySignature(query)]
	if !ok {
		return pgconn.NewCommandTag(""), errors.New("unexpected exec: " + querySignature(query))
	}
	return handler(args...)
}

type testRow struct {
	values []any
	err    error
}

func (r *testRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		if err := assignScanValue(dest[i], r.values[i]); err != nil {
			return err
		}
	}
	return nil
}

type testRows struct {
	rows   [][]any
	index  int
	closed bool
	err    error
}

func (r *testRows) Close() { r.closed = true }
func (r *testRows) Err() error {
	return r.err
}
func (r *testRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}
func (r *testRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *testRows) Next() bool {
	if r.err != nil || r.index >= len(r.rows) {
		r.closed = true
		return false
	}
	r.index++
	return true
}
func (r *testRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.rows) {
		return errors.New("scan called without current row")
	}
	current := r.rows[r.index-1]
	if len(dest) != len(current) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		if err := assignScanValue(dest[i], current[i]); err != nil {
			return err
		}
	}
	return nil
}
func (r *testRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.rows) {
		return nil, errors.New("values called without current row")
	}
	return r.rows[r.index-1], nil
}
func (r *testRows) RawValues() [][]byte { return nil }
func (r *testRows) Conn() *pgx.Conn     { return nil }

func assignScanValue(dest any, value any) error {
	destValue := reflect.ValueOf(dest)
	if destValue.Kind() != reflect.Ptr || destValue.IsNil() {
		return errors.New("scan destination must be a non-nil pointer")
	}
	target := destValue.Elem()
	if value == nil {
		target.SetZero()
		return nil
	}
	source := reflect.ValueOf(value)
	if source.Type().AssignableTo(target.Type()) {
		target.Set(source)
		return nil
	}
	if source.Type().ConvertibleTo(target.Type()) {
		target.Set(source.Convert(target.Type()))
		return nil
	}
	return errors.New("cannot assign scan value")
}

func newTestService(t *testing.T, dbtx db.DBTX, cfg *config.Config) *service.Service {
	t.Helper()

	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e3,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	require.NoError(t, err)
	t.Cleanup(cache.Close)

	svc := &service.Service{
		Config: cfg,
		Cache:  cache,
	}
	setUnexportedField(t, svc, "queries", db.New(dbtx))
	return svc
}

func setUnexportedField(t *testing.T, target any, fieldName string, value any) {
	t.Helper()

	v := reflect.ValueOf(target).Elem().FieldByName(fieldName)
	require.True(t, v.IsValid(), "field %s must exist", fieldName)
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}

func newTestConfig() *config.Config {
	return &config.Config{
		Env:                        "dev",
		APIProtocol:                "https",
		APIHost:                    "api.example.com",
		APIPort:                    8443,
		UIProtocol:                 "https",
		UIHost:                     "ui.example.com",
		UIPort:                     3000,
		APICORSMaxAgeSec:           300,
		APIRateLimitRequests:       100,
		APIRateLimitWindowSec:      60,
		APIThrottleConcurrentLimit: 10,
		APIRequestSizeLimitBytes:   1024 * 1024,
		APICacheDefaultCost:        1024,
		APICacheSourceTTLSec:       300,
		APICacheEventTTLSec:        900,
	}
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeContext))
}

func decodeJSONResponse[T any](t *testing.T, recorder *httptest.ResponseRecorder) T {
	t.Helper()

	var out T
	require.NoError(t, json.Unmarshal(recorder.Body.Bytes(), &out))
	return out
}

func querySignature(query string) string {
	firstLine, _, _ := strings.Cut(query, "\n")
	return firstLine
}

func testTime() time.Time {
	return time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
}
