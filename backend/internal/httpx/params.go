package httpx

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/vbsontakke/boilerplate/internal/apperror"
)

// ---- path params ---------------------------------------------------------

// PathUUID reads a chi URL param (e.g. /items/{id}) and parses it as a UUID.
// An invalid UUID returns 400 (use 404 instead if you prefer to hide it).
func PathUUID(r *http.Request, name string) (uuid.UUID, error) {
	raw := chi.URLParam(r, name)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperror.Validation("invalid path parameter",
			apperror.FieldError{Field: name, Issue: "must be a valid UUID"})
	}
	return id, nil
}

// ---- query params --------------------------------------------------------

// QueryString returns a trimmed query param, or "" if absent.
func QueryString(r *http.Request, key string) string {
	return strings.TrimSpace(r.URL.Query().Get(key))
}

// QueryInt64 returns (value, present, error) for an optional integer param.
func QueryInt64(r *http.Request, key string) (int64, bool, error) {
	raw := QueryString(r, key)
	if raw == "" {
		return 0, false, nil
	}
	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0, true, apperror.Validation("invalid query parameter",
			apperror.FieldError{Field: key, Issue: "must be an integer"})
	}
	return v, true, nil
}

// QueryEnum validates an optional param against allowed values.
func QueryEnum(r *http.Request, key string, allowed ...string) (string, error) {
	raw := QueryString(r, key)
	if raw == "" {
		return "", nil
	}
	for _, a := range allowed {
		if raw == a {
			return raw, nil
		}
	}
	return "", apperror.Validation("invalid query parameter",
		apperror.FieldError{Field: key, Issue: "must be one of: " + strings.Join(allowed, ", ")})
}

// ---- pagination ----------------------------------------------------------

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

// Pagination holds parsed ?limit=&cursor= values.
type Pagination struct {
	Limit  int
	Cursor string // opaque; decode with DecodeCursor
}

// ParsePagination reads ?limit (default 20, max 100) and ?cursor.
func ParsePagination(r *http.Request) (Pagination, error) {
	p := Pagination{Limit: DefaultLimit, Cursor: QueryString(r, "cursor")}
	if raw := QueryString(r, "limit"); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n < 1 || n > MaxLimit {
			return p, apperror.Validation("invalid query parameter",
				apperror.FieldError{Field: "limit", Issue: "must be an integer between 1 and 100"})
		}
		p.Limit = n
	}
	return p, nil
}

// EncodeCursor turns any keyset position (e.g. struct{CreatedAt; ID}) into an
// opaque URL-safe string for "next_cursor".
func EncodeCursor(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}

// DecodeCursor reverses EncodeCursor. A tampered cursor returns 400.
func DecodeCursor(cursor string, dst any) error {
	b, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil || json.Unmarshal(b, dst) != nil {
		return apperror.Validation("invalid query parameter",
			apperror.FieldError{Field: "cursor", Issue: "is malformed"})
	}
	return nil
}
