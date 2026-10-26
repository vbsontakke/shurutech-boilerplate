package httpx

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/vbsontakke/boilerplate/internal/apperror"
)

// The caller is identified by the X-User-ID header for now. Handlers only read
// the id through UserID(ctx), so the header can be replaced by JWT later
// without touching feature code.

const UserIDHeader = "X-User-ID"

type ctxKey int

const userIDKey ctxKey = iota

// RequireUserID is middleware that rejects requests without a valid
// X-User-ID header (401) and stores the id in the request context.
func RequireUserID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(strings.TrimSpace(r.Header.Get(UserIDHeader)))
		if err != nil {
			Error(w, r, apperror.ErrUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userIDKey, id)))
	})
}

// UserID returns the caller id set by RequireUserID.
func UserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

// OptionalUserID parses the header without requiring it (public endpoints
// that behave differently for the owner, e.g. viewing a draft).
func OptionalUserID(r *http.Request) (uuid.UUID, bool) {
	if id, ok := UserID(r.Context()); ok {
		return id, true
	}
	id, err := uuid.Parse(strings.TrimSpace(r.Header.Get(UserIDHeader)))
	return id, err == nil
}
