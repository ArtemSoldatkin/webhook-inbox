package service

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"
	"unsafe"

	"github.com/ArtemSoldatkin/webhook-inbox/internal/config"
	"github.com/ArtemSoldatkin/webhook-inbox/internal/db"
	"github.com/dgraph-io/ristretto"
	"github.com/go-openapi/testify/v2/assert"
	"github.com/go-openapi/testify/v2/require"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestNewService(t *testing.T) {
	t.Parallel()

	cfg := testServiceConfig()
	cache := testCache(t)

	svc := NewService(nil, cfg, cache)

	assert.Nil(t, svc.dbPool)
	assert.NotNil(t, svc.queries)
	assert.NotNil(t, svc.beginTx)
	assert.Same(t, cfg, svc.Config)
	assert.Same(t, cache, svc.Cache)
}

type serviceTestDB struct {
	queryHandlers    map[string]func(args ...any) (pgx.Rows, error)
	queryRowHandlers map[string]func(args ...any) pgx.Row
	execHandlers     map[string]func(args ...any) (pgconn.CommandTag, error)
}

func newServiceTestDB() *serviceTestDB {
	return &serviceTestDB{
		queryHandlers:    make(map[string]func(args ...any) (pgx.Rows, error)),
		queryRowHandlers: make(map[string]func(args ...any) pgx.Row),
		execHandlers:     make(map[string]func(args ...any) (pgconn.CommandTag, error)),
	}
}

func (dbtx *serviceTestDB) Query(_ context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	handler, ok := dbtx.queryHandlers[serviceQuerySignature(query)]
	if !ok {
		return nil, errors.New("unexpected query: " + serviceQuerySignature(query))
	}
	return handler(args...)
}

func (dbtx *serviceTestDB) QueryRow(_ context.Context, query string, args ...interface{}) pgx.Row {
	handler, ok := dbtx.queryRowHandlers[serviceQuerySignature(query)]
	if !ok {
		return &serviceTestRow{err: errors.New("unexpected query row: " + serviceQuerySignature(query))}
	}
	return handler(args...)
}

func (dbtx *serviceTestDB) Exec(_ context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	handler, ok := dbtx.execHandlers[serviceQuerySignature(query)]
	if !ok {
		return pgconn.NewCommandTag(""), errors.New("unexpected exec: " + serviceQuerySignature(query))
	}
	return handler(args...)
}

type serviceTestRow struct {
	values []any
	err    error
}

func (r *serviceTestRow) Scan(dest ...any) error {
	if r.err != nil {
		return r.err
	}
	if len(dest) != len(r.values) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		if err := serviceAssignScanValue(dest[i], r.values[i]); err != nil {
			return err
		}
	}
	return nil
}

type serviceTestRows struct {
	rows  [][]any
	index int
	err   error
}

func (r *serviceTestRows) Close()     {}
func (r *serviceTestRows) Err() error { return r.err }
func (r *serviceTestRows) CommandTag() pgconn.CommandTag {
	return pgconn.NewCommandTag("SELECT 0")
}
func (r *serviceTestRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (r *serviceTestRows) Next() bool {
	if r.err != nil || r.index >= len(r.rows) {
		return false
	}
	r.index++
	return true
}
func (r *serviceTestRows) Scan(dest ...any) error {
	if r.index == 0 || r.index > len(r.rows) {
		return errors.New("scan called without current row")
	}
	current := r.rows[r.index-1]
	if len(dest) != len(current) {
		return errors.New("scan destination count mismatch")
	}
	for i := range dest {
		if err := serviceAssignScanValue(dest[i], current[i]); err != nil {
			return err
		}
	}
	return nil
}
func (r *serviceTestRows) Values() ([]any, error) {
	if r.index == 0 || r.index > len(r.rows) {
		return nil, errors.New("values called without current row")
	}
	return r.rows[r.index-1], nil
}
func (r *serviceTestRows) RawValues() [][]byte { return nil }
func (r *serviceTestRows) Conn() *pgx.Conn     { return nil }

type serviceTestTx struct {
	*serviceTestDB
	commitCount   int
	rollbackCount int
}

func newServiceTestTx() *serviceTestTx {
	return &serviceTestTx{serviceTestDB: newServiceTestDB()}
}

func (tx *serviceTestTx) Begin(context.Context) (pgx.Tx, error) {
	return nil, errors.New("not implemented")
}
func (tx *serviceTestTx) Commit(context.Context) error {
	tx.commitCount++
	return nil
}
func (tx *serviceTestTx) Rollback(context.Context) error {
	tx.rollbackCount++
	return nil
}
func (tx *serviceTestTx) CopyFrom(context.Context, pgx.Identifier, []string, pgx.CopyFromSource) (int64, error) {
	return 0, errors.New("not implemented")
}
func (tx *serviceTestTx) SendBatch(context.Context, *pgx.Batch) pgx.BatchResults { return nil }
func (tx *serviceTestTx) LargeObjects() pgx.LargeObjects                         { return pgx.LargeObjects{} }
func (tx *serviceTestTx) Prepare(context.Context, string, string) (*pgconn.StatementDescription, error) {
	return nil, errors.New("not implemented")
}
func (tx *serviceTestTx) Conn() *pgx.Conn { return nil }

func serviceAssignScanValue(dest any, value any) error {
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

func newServiceUnderTest(t *testing.T, dbtx db.DBTX) *Service {
	t.Helper()

	cache := testCache(t)
	svc := &Service{
		Config: testServiceConfig(),
		Cache:  cache,
		beginTx: func(context.Context, pgx.TxOptions) (pgx.Tx, error) {
			return nil, errors.New("unexpected begin tx")
		},
	}
	setServicePrivateField(t, svc, "queries", db.New(dbtx))
	return svc
}

func testCache(t *testing.T) *ristretto.Cache {
	t.Helper()
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e3,
		MaxCost:     1 << 20,
		BufferItems: 64,
	})
	require.NoError(t, err)
	t.Cleanup(cache.Close)
	return cache
}

func testServiceConfig() *config.Config {
	return &config.Config{
		Env:                            "dev",
		APICacheDefaultCost:            1024,
		APICacheSourceTTLSec:           300,
		APICacheEventTTLSec:            900,
		APIDeliveryRetryBackoffBaseSec: 1,
		APIDeliveryRetryBackoffMaxSec:  60,
	}
}

func setServicePrivateField(t *testing.T, target any, fieldName string, value any) {
	t.Helper()
	v := reflect.ValueOf(target).Elem().FieldByName(fieldName)
	require.True(t, v.IsValid(), "field %s must exist", fieldName)
	reflect.NewAt(v.Type(), unsafe.Pointer(v.UnsafeAddr())).Elem().Set(reflect.ValueOf(value))
}

func serviceQuerySignature(query string) string {
	firstLine, _, _ := strings.Cut(query, "\n")
	return firstLine
}

func serviceSourceRowValues(source db.Source) []any {
	return []any{
		source.ID,
		source.PublicID,
		source.EgressUrl,
		source.StaticHeaders,
		source.Status,
		source.StatusReason,
		source.Description,
		source.CreatedAt,
		source.UpdatedAt,
		source.DisableAt,
	}
}

func serviceEventRowValues(event db.Event) []any {
	return []any{
		event.ID,
		event.SourceID,
		event.DedupHash,
		event.Method,
		event.IngressPath,
		event.RemoteAddress,
		event.QueryParams,
		event.RawHeaders,
		event.Body,
		event.BodyContentType,
		event.ReceivedAt,
	}
}

func mustServiceUUID(t *testing.T, value string) pgtype.UUID {
	t.Helper()
	var id pgtype.UUID
	require.NoError(t, id.Scan(value))
	return id
}

func testNow() time.Time {
	return time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)
}
