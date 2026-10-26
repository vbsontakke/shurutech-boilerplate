package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/vbsontakke/boilerplate/internal/apperror"
)

const maxBodyBytes = 1 << 20 // 1 MB

// DecodeJSON strictly decodes a JSON request body into dst:
// size-limited, unknown fields rejected, exactly one JSON object.
// Every failure is returned as a 400 validation error with a clear message.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		var maxErr *http.MaxBytesError

		switch {
		case errors.Is(err, io.EOF):
			return apperror.Validation("request body must not be empty")
		case errors.As(err, &syntaxErr):
			return apperror.Validation(fmt.Sprintf("malformed JSON at position %d", syntaxErr.Offset))
		case errors.Is(err, io.ErrUnexpectedEOF):
			return apperror.Validation("malformed JSON")
		case errors.As(err, &typeErr):
			return apperror.Validation("invalid value type",
				apperror.FieldError{Field: typeErr.Field, Issue: "must be " + typeErr.Type.String()})
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.Trim(strings.TrimPrefix(err.Error(), "json: unknown field "), `"`)
			return apperror.Validation("unknown field in request body",
				apperror.FieldError{Field: field, Issue: "unknown field"})
		case errors.As(err, &maxErr):
			return apperror.Validation("request body too large")
		default:
			return apperror.Validation("invalid request body")
		}
	}

	// Anything after the first value (even a second object) is rejected.
	var extra json.RawMessage
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		return apperror.Validation("request body must contain a single JSON object")
	}
	return nil
}

// DecodeAndValidate decodes the body and runs struct-tag validation.
func DecodeAndValidate(w http.ResponseWriter, r *http.Request, dst any) error {
	if err := DecodeJSON(w, r, dst); err != nil {
		return err
	}
	return Validate(dst)
}
