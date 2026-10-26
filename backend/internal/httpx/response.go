// Package httpx holds reusable HTTP plumbing: JSON responses, a single
// error→status mapper, request decoding/validation, query helpers and
// middleware. Feature handlers should only ever respond through these.
package httpx

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

// ErrorBody is the one error shape every endpoint returns:
//
//	{"error":{"code":"NOT_FOUND","message":"ad not found","details":[...]}}
type ErrorBody struct {
	Error ErrorPayload `json:"error"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

// Page is the standard list envelope. next_cursor is "" on the last page:
//
//	{"data":[...],"next_cursor":"...","limit":20}
type Page[T any] struct {
	Data       []T    `json:"data"`
	NextCursor string `json:"next_cursor"`
	Limit      int    `json:"limit"`
}

// NewPage guarantees "data" is [] rather than null when empty.
func NewPage[T any](items []T, nextCursor string, limit int) Page[T] {
	if items == nil {
		items = []T{}
	}
	return Page[T]{Data: items, NextCursor: nextCursor, Limit: limit}
}

// JSON writes v as JSON with the given status code.
func JSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if v == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("encode response", "err", err)
	}
}

// OK / Created / NoContent are small conveniences for readability.
func OK(w http.ResponseWriter, v any)      { JSON(w, http.StatusOK, v) }
func Created(w http.ResponseWriter, v any) { JSON(w, http.StatusCreated, v) }
func NoContent(w http.ResponseWriter)      { w.WriteHeader(http.StatusNoContent) }
