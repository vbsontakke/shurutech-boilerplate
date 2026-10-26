# Prompt to regenerate this boilerplate with an AI coding agent

Use this if you want to rebuild the boilerplate from scratch (e.g. in a different
language/framework the interviewer specifies). It contains **no problem-specific
design**, only generic infrastructure, which the guidelines allow you to prepare.

---

> Create a production-style **generic backend boilerplate** (no business features) in a
> monorepo with `backend/` (Go 1.24) and `frontend/` (React + TypeScript via Vite).
>
> **Backend stack:** chi v5 router, go-chi/cors, pgx v5 (`pgxpool`) for PostgreSQL 16,
> godotenv, go-playground/validator v10, google/uuid, log/slog, testify. No ORM.
>
> **Requirements:**
> 1. `cmd/server/main.go`: load config → slog logger (text in development, JSON otherwise) →
>    pgx pool with ping → run embedded SQL migrations → chi router → `http.Server` with
>    read/write/idle timeouts and graceful shutdown on SIGINT/SIGTERM.
> 2. `internal/config`: read APP_ENV, APP_PORT, LOG_LEVEL, CORS_ORIGINS, DB_HOST, DB_PORT,
>    DB_USER, DB_PASSWORD, DB_NAME, DB_SSLMODE, DB_MAX_CONNS from env (+ optional `.env`).
>    Fail fast if DB_USER/DB_PASSWORD/DB_NAME are missing. Build the DSN with `net/url` so
>    credentials are escaped. No secrets in code; provide `.env.example`.
> 3. `internal/database`: `NewPool`, a `DBTX` interface satisfied by pool and `pgx.Tx`,
>    `WithTx(ctx, pool, fn)` that commits/rolls back (also on panic), a tiny migration runner
>    applying `migrations/*.up.sql` in order inside transactions, tracked in
>    `schema_migrations`, guarded by `pg_advisory_lock`. Helpers: `IsNoRows`,
>    `IsUniqueViolation`, `IsForeignKeyViolation`, `IsCheckViolation`, `ConstraintName`.
> 4. `migrations/`: `embed.FS`; `0001_init.up.sql` enabling pgcrypto and a reusable
>    `set_updated_at()` trigger function.
> 5. `internal/apperror`: transport-agnostic `Error{Kind, Code, Message, Fields, Err}` with
>    kinds Internal/Validation/Unauthorized/Forbidden/NotFound/Conflict/Unprocessable/
>    TooManyRequests, constructors, `Wrap`, `From(err)`, and `Is` matching by Kind+Code so
>    sentinels work with `errors.Is`.
> 6. `internal/httpx`: `JSON/OK/Created/NoContent`; a single `StatusFor(kind)` mapper; `Error(w,r,err)`
>    writing `{"error":{"code","message","details"}}` and hiding internal errors (logged with
>    request id); JSON 404/405 handlers; strict `DecodeJSON` (1MB limit, unknown fields rejected,
>    single object, friendly messages) and `DecodeAndValidate` using validator with JSON field
>    names and readable issues; `PathUUID`, `QueryString`, `QueryInt64`, `QueryEnum`;
>    `ParsePagination` (limit default 20, max 100, opaque cursor) with `EncodeCursor/DecodeCursor`;
>    generic `Page[T]{data,next_cursor,limit}`; middleware: structured request logger and a
>    panic recoverer returning JSON 500; `RequireUserID` middleware reading `X-User-ID`.
> 7. `internal/server`: router with RequestID, RealIP, logger, recoverer, 30s timeout, CORS;
>    `GET /health` → `{"status":"running"}`, `GET /health/db` → DB ping (503 if down);
>    an empty `/api/v1` group with comments showing how to register a feature.
> 8. Tests: handler tests via httptest for /health, 404, 405; table-driven tests for decode/
>    validation, status mapping, pagination, cursor round trip, RequireUserID; DB integration
>    tests (connect+migrate, WithTx rollback, unique violation) using a `testutil.NewTestDB`
>    helper that **skips** when the test DB is unreachable. Also `testutil.Do` and
>    generic `testutil.DecodeJSON[T]`.
> 9. Tooling: root `docker-compose.yml` with `postgres` (5432, volume) and `postgres-test`
>    (5433, tmpfs) with healthchecks; `backend/Makefile` (setup, db-up, db-down, db-reset,
>    db-psql, run, dev with air, build, test, test-v, cover, lint, fmt, new-migration);
>    `.air.toml`; `api.http`; `.gitignore`.
>
> **Frontend:** Vite React TS app; Vite dev proxy `/api` and `/health` → `http://localhost:8080`;
> `src/api/client.ts` fetch wrapper that throws a typed `ApiError` from the backend error shape,
> supports query params and an optional `X-User-ID` header; `App.tsx` showing API and DB health
> with a re-check button; minimal CSS with dark mode. Remove Vite demo assets.
>
> **Docs:** `README.md` (quickstart, commands, layout, how to add a feature), `CLAUDE.md`
> (conventions for AI agents: layering, error handling, SQL rules, testing), and an empty
> `DESIGN.md` template (assumptions, entities, state machine, numbered rules, API table,
> schema, concurrency, build order).
>
> Verify: `go build ./...`, `go vet ./...`, `go test ./...` pass (with and without the test DB),
> `npm run build` passes, and `curl localhost:8080/health` returns `{"status":"running"}`.
