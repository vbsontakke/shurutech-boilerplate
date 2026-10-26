package httpx

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"

	"github.com/vbsontakke/boilerplate/internal/apperror"
)

// validate is safe for concurrent use and caches struct metadata.
var validate = newValidator()

func newValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	// Report JSON field names ("min_price") instead of Go names ("MinPrice").
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name := strings.SplitN(f.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return f.Name
		}
		return name
	})
	return v
}

// Validate runs `validate:"..."` struct tags and converts failures into a
// 400 apperror with one FieldError per invalid field.
//
// Example DTO:
//
//	type CreateThingRequest struct {
//	    Title string  `json:"title" validate:"required,min=5,max=100"`
//	    Price int64   `json:"price" validate:"required,gt=0"`
//	    Email string  `json:"email" validate:"required,email"`
//	    Tags  []string `json:"tags"  validate:"max=5,dive,required"`
//	}
func Validate(v any) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return apperror.Internal(err)
	}

	fields := make([]apperror.FieldError, 0, len(verrs))
	for _, fe := range verrs {
		fields = append(fields, apperror.FieldError{Field: fieldPath(fe), Issue: issueFor(fe)})
	}
	return apperror.Validation("request validation failed", fields...)
}

// fieldPath drops the top-level struct name: "CreateReq.title" → "title".
func fieldPath(fe validator.FieldError) string {
	ns := fe.Namespace()
	if i := strings.Index(ns, "."); i >= 0 {
		return ns[i+1:]
	}
	return fe.Field()
}

func issueFor(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "uuid", "uuid4":
		return "must be a valid UUID"
	case "min":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must be at least %s characters", fe.Param())
		}
		return "must be at least " + fe.Param()
	case "max":
		if fe.Kind() == reflect.String {
			return fmt.Sprintf("must be at most %s characters", fe.Param())
		}
		return "must be at most " + fe.Param()
	case "gt":
		return "must be greater than " + fe.Param()
	case "gte":
		return "must be greater than or equal to " + fe.Param()
	case "lt":
		return "must be less than " + fe.Param()
	case "lte":
		return "must be less than or equal to " + fe.Param()
	case "oneof":
		return "must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	case "len":
		return "must have length " + fe.Param()
	default:
		return "failed '" + fe.Tag() + "' validation"
	}
}
