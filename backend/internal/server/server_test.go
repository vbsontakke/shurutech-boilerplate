package server_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vbsontakke/boilerplate/internal/httpx"
	"github.com/vbsontakke/boilerplate/internal/server"
	"github.com/vbsontakke/boilerplate/internal/testutil"
)

// Sample API tests: copy this pattern for feature endpoints.

func newTestRouter() http.Handler {
	return server.NewRouter(server.Deps{Log: testutil.DiscardLogger(), DB: nil})
}

func TestHealth_ReturnsRunning(t *testing.T) {
	rec := testutil.Do(t, newTestRouter(), testutil.Request{Method: http.MethodGet, Path: "/health"})

	require.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Header().Get("Content-Type"), "application/json")
	assert.JSONEq(t, `{"status":"running"}`, rec.Body.String())
}

func TestHealthDB_WithoutDatabase_Returns503(t *testing.T) {
	rec := testutil.Do(t, newTestRouter(), testutil.Request{Method: http.MethodGet, Path: "/health/db"})
	assert.Equal(t, http.StatusServiceUnavailable, rec.Code)
}

func TestUnknownRoute_ReturnsJSON404(t *testing.T) {
	rec := testutil.Do(t, newTestRouter(), testutil.Request{Method: http.MethodGet, Path: "/does-not-exist"})

	require.Equal(t, http.StatusNotFound, rec.Code)
	body := testutil.DecodeJSON[httpx.ErrorBody](t, rec)
	assert.Equal(t, "NOT_FOUND", body.Error.Code)
}

func TestWrongMethod_ReturnsJSON405(t *testing.T) {
	rec := testutil.Do(t, newTestRouter(), testutil.Request{Method: http.MethodPost, Path: "/health"})

	require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	body := testutil.DecodeJSON[httpx.ErrorBody](t, rec)
	assert.Equal(t, "METHOD_NOT_ALLOWED", body.Error.Code)
}
