// Package health exposes liveness and readiness endpoints.
package health

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/vbsontakke/boilerplate/internal/httpx"
)

// Pinger is satisfied by *pgxpool.Pool; an interface keeps this testable.
type Pinger interface {
	Ping(ctx context.Context) error
}

type Handler struct {
	db Pinger
}

func NewHandler(db Pinger) *Handler { return &Handler{db: db} }

func (h *Handler) Routes(r chi.Router) {
	r.Get("/health", h.Live)
	r.Get("/health/db", h.Ready)
}

// Live: GET /health → 200 {"status":"running"} (process is up).
func (h *Handler) Live(w http.ResponseWriter, _ *http.Request) {
	httpx.OK(w, map[string]string{"status": "running"})
}

// Ready: GET /health/db → 200 if the database answers, else 503.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable", "database": "not configured"})
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()
	if err := h.db.Ping(ctx); err != nil {
		httpx.JSON(w, http.StatusServiceUnavailable, map[string]string{"status": "unavailable", "database": "down"})
		return
	}
	httpx.OK(w, map[string]string{"status": "running", "database": "up"})
}
