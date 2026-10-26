package httpx

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/vbsontakke/boilerplate/internal/apperror"
)

// StatusFor is the ONE place where application error kinds become HTTP
// status codes. Services never choose status codes themselves.
func StatusFor(k apperror.Kind) int {
	switch k {
	case apperror.KindValidation:
		return http.StatusBadRequest
	case apperror.KindUnauthorized:
		return http.StatusUnauthorized
	case apperror.KindForbidden:
		return http.StatusForbidden
	case apperror.KindNotFound:
		return http.StatusNotFound
	case apperror.KindConflict:
		return http.StatusConflict
	case apperror.KindUnprocessable:
		return http.StatusUnprocessableEntity
	case apperror.KindTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

// Error writes err using the standard error body. Internal errors are
// logged with the request id and returned to the client as a generic message.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	ae := apperror.From(err)
	status := StatusFor(ae.Kind)

	if status >= 500 {
		slog.ErrorContext(r.Context(), "request failed",
			"err", err, "path", r.URL.Path, "method", r.Method,
			"request_id", middleware.GetReqID(r.Context()))
	}

	payload := ErrorPayload{Code: ae.Code, Message: ae.Message}
	if len(ae.Fields) > 0 {
		payload.Details = ae.Fields
	}
	JSON(w, status, ErrorBody{Error: payload})
}

// NotFoundHandler / MethodNotAllowedHandler keep router-level errors in the
// same JSON shape as everything else.
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	Error(w, r, apperror.NotFound("route not found"))
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusMethodNotAllowed, ErrorBody{Error: ErrorPayload{
		Code: "METHOD_NOT_ALLOWED", Message: "method not allowed for this route",
	}})
}
