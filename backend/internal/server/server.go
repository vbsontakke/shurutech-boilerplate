// Package server wires middleware and routes into one http.Handler.
package server

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/vbsontakke/boilerplate/internal/health"
	"github.com/vbsontakke/boilerplate/internal/httpx"
)

// Deps are the shared dependencies handlers may need. Add services here as
// features are built (constructed in cmd/server/main.go).
type Deps struct {
	Log         *slog.Logger
	DB          *pgxpool.Pool // may be nil in unit tests
	CORSOrigins []string
}

// NewRouter builds the HTTP handler with global middleware and routes.
func NewRouter(d Deps) http.Handler {
	r := chi.NewRouter()

	// ---- global middleware (order matters) ----
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(httpx.RequestLogger(d.Log))
	r.Use(httpx.Recoverer(d.Log))
	r.Use(middleware.Timeout(30 * time.Second))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   d.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-User-ID", "Idempotency-Key"},
		ExposedHeaders:   []string{"X-Request-Id", "Location"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Consistent JSON for unknown routes / wrong methods.
	r.NotFound(httpx.NotFoundHandler)
	r.MethodNotAllowed(httpx.MethodNotAllowedHandler)

	// ---- routes ----
	var pinger health.Pinger
	if d.DB != nil {
		pinger = d.DB
	}
	health.NewHandler(pinger).Routes(r)

	r.Route("/api/v1", func(api chi.Router) {
		// Register feature routes here, e.g.:
		//
		//   itemRepo := item.NewRepository(d.DB)
		//   itemSvc  := item.NewService(itemRepo)
		//   item.NewHandler(itemSvc).Routes(api)
		//
		// Routes that need the caller:
		//   api.Group(func(pr chi.Router) {
		//       pr.Use(httpx.RequireUserID)
		//       ...
		//   })
	})

	return r
}
