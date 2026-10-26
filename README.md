# Backend Boilerplate: Go API + PostgreSQL + React

A generic, ready-to-extend starter. It contains **infrastructure only**: no business
features. Write your domain model, API contracts and schema in `DESIGN.md` first, then
build features on top of the conventions below.

| Layer | Stack |
|---|---|
| API | Go 1.24, chi v5, slog, go-playground/validator |
| DB | PostgreSQL 16, pgx v5 (no ORM), embedded SQL migrations |
| Tests | `testing` + testify, httptest, auto-skipping DB integration tests |
| Frontend | React + TypeScript (Vite), dev proxy to the API |
| Dev | Docker Compose (dev + test DB), Makefile, air hot reload |

---

## 1. Prerequisites
- Go **1.24+** · Docker Desktop · Node **20+** · VS Code (+ Go, REST Client, Claude Code extensions)

## 2. Quickstart

```bash
# backend
cd backend
make setup          # creates .env, runs go mod tidy (writes go.sum), installs air
make db-up          # starts Postgres :5432 (dev) and :5433 (test), waits until healthy
make run            # applies migrations, starts API on :8080

curl localhost:8080/health      # {"status":"running"}
curl localhost:8080/health/db   # {"database":"up","status":"running"}

make test           # unit + API + DB integration tests
```

```bash
# frontend (new terminal)
cd frontend
npm install
npm run dev         # http://localhost:5173 shows API + DB status
```

> No Docker? Point `DB_*` in `backend/.env` at any local Postgres, and `TEST_DB_*` at a
> separate test database. DB tests skip automatically if the test DB is unreachable.

## 3. Commands (`backend/`)

| Command | What it does |
|---|---|
| `make run` / `make dev` | Run API / run with hot reload |
| `make test` / `make test-v` | Tests / verbose + race detector |
| `make cover` | Coverage summary |
| `make lint` / `make fmt` | `go vet` + gofmt check / format |
| `make db-up` / `db-down` / `db-reset` | Start / stop / wipe databases |
| `make db-psql` | psql shell on dev DB |
| `make new-migration name=create_items` | Creates `000N_create_items.{up,down}.sql` |

## 4. Layout

```
├── CLAUDE.md                 conventions for AI agents (Claude Code reads this)
├── DESIGN.md                 design template, filled in before building features
├── PROMPT.md                 prompt to regenerate this boilerplate
├── docker-compose.yml        postgres (dev) + postgres-test
├── backend/
│   ├── cmd/server/main.go    config → logger → pool → migrations → router → graceful server
│   ├── internal/
│   │   ├── config/           env loading (.env optional), validated
│   │   ├── database/         pgx pool, WithTx, Migrate, pg error helpers
│   │   ├── apperror/         transport-agnostic errors (Kind, Code, Message, Fields)
│   │   ├── httpx/            JSON, StatusFor, Error, DecodeAndValidate, params,
│   │   │                     pagination/cursor, middleware, X-User-ID identity
│   │   ├── health/           GET /health, GET /health/db
│   │   ├── server/           router + middleware + route registration
│   │   └── testutil/         Do(), DecodeJSON[T](), NewTestDB(), Truncate()
│   ├── migrations/           0001_init.up.sql … (embedded, run on startup)
│   ├── api.http              clickable requests (REST Client)
│   └── Makefile
└── frontend/
    ├── vite.config.ts        proxy /api, /health → :8080
    └── src/api/client.ts     typed fetch wrapper (ApiError, Page<T>, X-User-ID)
```

## 5. Built-in conventions

**Error response (every endpoint)**
```json
{"error":{"code":"VALIDATION_FAILED","message":"request validation failed",
          "details":[{"field":"title","issue":"must be at least 5 characters"}]}}
```

**Status mapping (one place, `httpx.StatusFor`)**

| apperror | HTTP | apperror | HTTP |
|---|---|---|---|
| `Validation` | 400 | `Conflict` | 409 |
| `Unauthorized` | 401 | `Unprocessable` | 422 |
| `Forbidden` | 403 | `TooManyRequests` | 429 |
| `NotFound` | 404 | anything else | 500 (details hidden, logged) |

**List response:** `{"data":[...],"next_cursor":"...","limit":20}` via `httpx.NewPage`.

## 6. Adding a feature

1. `make new-migration name=create_things` and write tables, constraints and indexes.
2. Create `internal/thing/`:
   - `model.go`: structs, status constants, `CanTransition`, feature sentinel errors
   - `repository.go`: interface + pgx implementation (SQL only)
   - `service.go`: business rules, ownership, transactions (`database.WithTx`)
   - `handler.go`: DTOs with `validate` tags, `Routes(r chi.Router)`
3. Register in `internal/server/server.go` inside `/api/v1`.
4. Tests: service (fake repo) → handler (`testutil.Do`) → optional DB integration.

Handler skeleton:
```go
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
    userID, _ := httpx.UserID(r.Context())
    var req CreateThingRequest
    if err := httpx.DecodeAndValidate(w, r, &req); err != nil {
        httpx.Error(w, r, err); return
    }
    thing, err := h.svc.Create(r.Context(), userID, req.toInput())
    if err != nil { httpx.Error(w, r, err); return }
    httpx.Created(w, toResponse(thing))
}
```

## 7. Checklist before the first feature
- [ ] `make db-up && make run`, then `/health` returns `{"status":"running"}`
- [ ] `make test` is green
- [ ] Module path changed from `github.com/vbsontakke/boilerplate` (`go.mod` and imports) if this is not your repo
- [ ] `DESIGN.md` filled in before writing feature code
- [ ] Commit after each vertical slice
