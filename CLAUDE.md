# Project conventions (read before writing code)

Go REST API + PostgreSQL, optional React (Vite) frontend. The boilerplate is done;
features are added on top of it. The domain model, API contracts and schema come
from `DESIGN.md`: implement that design exactly; do not invent extra rules.

## Layout
```
backend/
  cmd/server/main.go        wiring only (config → pool → migrations → router → server)
  internal/apperror/        transport-agnostic errors (Kind + Code + Message)
  internal/httpx/           JSON responses, error→status mapper, decode/validate, params, middleware
  internal/database/        pgx pool, WithTx, Migrate, pg error helpers (IsUniqueViolation…)
  internal/server/          router: middleware + route registration
  internal/<feature>/       one package per feature (see below)
  migrations/               0001_init.up.sql, 0002_<name>.up.sql … (embedded, auto-applied)
frontend/src/api/client.ts  fetch wrapper matching the API's JSON contract
```

## Adding a feature `internal/<feature>/`
- `model.go`: domain structs, enums/status constants, state-machine helpers, feature sentinel
  errors (e.g. `var ErrInvalidTransition = apperror.Conflict("INVALID_TRANSITION", "...")`).
- `repository.go`: `Repository` interface + pgx implementation. SQL only, no business rules.
  Accept `database.DBTX` so methods work inside or outside a transaction.
- `service.go`: business rules, ownership checks, transitions, transactions (`database.WithTx`).
  Never imports `net/http`. Returns `*apperror.Error` (or wraps DB errors with `apperror.Internal`).
- `handler.go`: `Routes(r chi.Router)`; decode with `httpx.DecodeAndValidate`, path ids with
  `httpx.PathUUID`, pagination with `httpx.ParsePagination`, respond with `httpx.OK/Created/NoContent`,
  errors with `httpx.Error(w, r, err)`. Request/response DTOs live here, separate from domain structs.
- Register routes in `internal/server/server.go` under `/api/v1`. Caller-identified routes go
  inside a group using `httpx.RequireUserID`.
- New tables: `make new-migration name=<name>`, then write SQL with constraints (NOT NULL, CHECK,
  UNIQUE, FK) and indexes for the queries actually used. Use the `set_updated_at()` trigger.

## Rules
- Money: `BIGINT` minor units / `int64`. Never float.
- Timestamps: `TIMESTAMPTZ`, UTC.
- SQL: parameterised (`$1, $2`) only. Never build SQL from user input with string concat;
  dynamic filters append `$n` placeholders.
- Multi-statement writes run in one transaction. Race-prone rules use `SELECT … FOR UPDATE`
  or a conditional `UPDATE … WHERE <expected state>` and check `RowsAffected()`.
- PATCH DTOs use pointer fields to distinguish "absent" from "zero".
- Status codes come only from `httpx.StatusFor` via apperror kinds:
  400 validation · 401 no identity · 403 not allowed · 404 missing/hidden · 409 state conflict/duplicate ·
  422 business rule · 500 unexpected.
- Lists return `httpx.NewPage(items, nextCursor, limit)` → `{"data":[],"next_cursor":"","limit":20}`.
- Keep functions small; no global state beyond the validator.

## Tests
- Service: table-driven tests with a hand-written fake `Repository`.
- Handler: `testutil.Do(t, router, testutil.Request{...})` + `testutil.DecodeJSON[T]`.
- DB integration: `testutil.NewTestDB(t)` (auto-skips if test DB down) + `testutil.Truncate`.
- After every change run `make test` (from `backend/`) and `go vet ./...`.

## When you finish a step
List every file you created or changed, and call out any assumption not stated in DESIGN.md.
