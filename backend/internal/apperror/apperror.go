// Package apperror defines transport-agnostic application errors.
//
// Services return these; they never import net/http. The HTTP layer
// (internal/httpx) maps each Kind to a status code in exactly one place.
package apperror

import (
	"errors"
	"fmt"
)

type Kind int

const (
	KindInternal        Kind = iota // 500 – unexpected; message is never shown to clients
	KindValidation                  // 400 – malformed input / field validation failed
	KindUnauthorized                // 401 – caller identity missing or invalid
	KindForbidden                   // 403 – authenticated but not allowed
	KindNotFound                    // 404 – resource does not exist (or is hidden)
	KindConflict                    // 409 – state conflict / duplicate / invalid transition
	KindUnprocessable               // 422 – valid input that breaks a business rule
	KindTooManyRequests             // 429 – rate limited
)

// FieldError describes one invalid input field.
type FieldError struct {
	Field string `json:"field"`
	Issue string `json:"issue"`
}

// Error is the single error type services return.
type Error struct {
	Kind    Kind
	Code    string       // stable machine-readable code, e.g. "INVALID_TRANSITION"
	Message string       // human-readable, safe to show to clients
	Fields  []FieldError // optional per-field details (validation)
	Err     error        // wrapped cause, logged but never returned to clients
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.Err }

// Is lets errors.Is(err, apperror.ErrNotFound) match by Kind+Code,
// so sentinel values below can be compared even after wrapping.
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Kind == t.Kind && e.Code == t.Code
}

// Wrap attaches an underlying cause (kept for logs) and returns a copy.
func (e *Error) Wrap(cause error) *Error {
	cp := *e
	cp.Err = cause
	return &cp
}

// WithMessage returns a copy with a more specific client-facing message.
func (e *Error) WithMessage(msg string) *Error {
	cp := *e
	cp.Message = msg
	return &cp
}

// ---- constructors -------------------------------------------------------

func New(kind Kind, code, msg string) *Error { return &Error{Kind: kind, Code: code, Message: msg} }

func Validation(msg string, fields ...FieldError) *Error {
	return &Error{Kind: KindValidation, Code: "VALIDATION_FAILED", Message: msg, Fields: fields}
}
func Unauthorized(msg string) *Error   { return New(KindUnauthorized, "UNAUTHORIZED", msg) }
func Forbidden(msg string) *Error      { return New(KindForbidden, "FORBIDDEN", msg) }
func NotFound(msg string) *Error       { return New(KindNotFound, "NOT_FOUND", msg) }
func Conflict(code, msg string) *Error { return New(KindConflict, code, msg) }
func Unprocessable(code, msg string) *Error {
	return New(KindUnprocessable, code, msg)
}
func Internal(cause error) *Error {
	return &Error{Kind: KindInternal, Code: "INTERNAL_ERROR", Message: "something went wrong", Err: cause}
}

// From converts any error into *Error. Unknown errors become Internal.
func From(err error) *Error {
	if err == nil {
		return nil
	}
	var ae *Error
	if errors.As(err, &ae) {
		return ae
	}
	return Internal(err)
}

// ---- common sentinels (compare with errors.Is) ---------------------------
// Add domain-specific ones next to your feature, e.g.
//   var ErrInvalidTransition = apperror.Conflict("INVALID_TRANSITION", "...")

var (
	ErrNotFound     = NotFound("resource not found")
	ErrForbidden    = Forbidden("you are not allowed to perform this action")
	ErrUnauthorized = Unauthorized("missing or invalid caller identity")
)
