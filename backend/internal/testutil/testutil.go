// Package testutil provides helpers for API and database tests.
package testutil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vbsontakke/boilerplate/internal/config"
	"github.com/vbsontakke/boilerplate/internal/database"
	"github.com/vbsontakke/boilerplate/migrations"
)

// DiscardLogger keeps test output clean.
func DiscardLogger() *slog.Logger { return slog.New(slog.NewTextHandler(io.Discard, nil)) }

// ---- HTTP helpers --------------------------------------------------------

// Request is a small builder for handler tests.
type Request struct {
	Method  string
	Path    string
	Body    any               // marshalled to JSON if not nil; string sent raw
	Headers map[string]string // e.g. {"X-User-ID": "..."}
}

// Do executes req against h and returns the recorder.
func Do(t *testing.T, h http.Handler, req Request) *httptest.ResponseRecorder {
	t.Helper()
	var body io.Reader
	switch b := req.Body.(type) {
	case nil:
	case string:
		body = strings.NewReader(b)
	default:
		raw, err := json.Marshal(b)
		if err != nil {
			t.Fatalf("marshal body: %v", err)
		}
		body = bytes.NewReader(raw)
	}
	r := httptest.NewRequest(req.Method, req.Path, body)
	if body != nil {
		r.Header.Set("Content-Type", "application/json")
	}
	for k, v := range req.Headers {
		r.Header.Set(k, v)
	}
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, r)
	return rec
}

// DecodeJSON decodes a recorder body into T.
func DecodeJSON[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("decode response %q: %v", rec.Body.String(), err)
	}
	return v
}

// ---- database helpers ----------------------------------------------------

// NewTestDB connects to the TEST database (TEST_DB_* env vars, defaults
// match docker-compose's postgres-test on :5433), runs migrations and
// returns a pool. The test is SKIPPED if the database is unreachable, so
// `go test ./...` still passes without Docker.
func NewTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
	cfg := config.DBConfig{
		Host:     envOr("TEST_DB_HOST", "localhost"),
		Port:     envOr("TEST_DB_PORT", "5433"),
		User:     envOr("TEST_DB_USER", "app"),
		Password: envOr("TEST_DB_PASSWORD", "app"),
		Name:     envOr("TEST_DB_NAME", "app_test"),
		SSLMode:  "disable",
		MaxConns: 5,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := database.NewPool(ctx, cfg)
	if err != nil {
		t.Skipf("test database unavailable (%v); run `make db-up` to enable integration tests", err)
	}
	if err := database.Migrate(ctx, pool, migrations.FS, nil); err != nil {
		pool.Close()
		t.Fatalf("migrate test db: %v", err)
	}
	t.Cleanup(pool.Close)
	return pool
}

// Truncate empties the given tables (and resets identities) between tests.
//
//	testutil.Truncate(t, pool, "orders", "users")
func Truncate(t *testing.T, pool *pgxpool.Pool, tables ...string) {
	t.Helper()
	if len(tables) == 0 {
		return
	}
	quoted := make([]string, len(tables))
	for i, tbl := range tables {
		quoted[i] = pgx.Identifier{tbl}.Sanitize()
	}
	sql := "TRUNCATE " + strings.Join(quoted, ", ") + " RESTART IDENTITY CASCADE"
	if _, err := pool.Exec(context.Background(), sql); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}

func envOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
