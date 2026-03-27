package middleware

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"reflect"
	"strings"

	commonapi "go-gin-ecommerce/internal/common/api"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		if c.Writer.Written() && c.Writer.Size() > 0 {
			return
		}

		apiErr := normalizeError(c.Errors.Last().Err)

		if apiErr.Status >= 500 {
			logger.Error("request failed", "path", c.Request.URL.Path, "status", apiErr.Status, "error", c.Errors.Last().Err)
		}

		c.AbortWithStatusJSON(apiErr.Status, commonapi.NewErrorResponse(c.Request.URL.Path, apiErr))
	}
}

func normalizeError(err error) *commonapi.Error {
	if err == nil {
		return commonapi.NewInternalServerError()
	}

	var apiErr *commonapi.Error
	if errors.As(err, &apiErr) {
		return apiErr
	}

	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return commonapi.NewValidationError(validationDetails(validationErrs))
	}

	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "body",
			Constraints: []string{"body must contain valid JSON"},
		}})
	}

	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		field := typeErr.Field
		if field == "" {
			field = "body"
		}

		return commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       field,
			Constraints: []string{fmt.Sprintf("%s must be %s", field, humanizeType(typeErr.Type))},
		}})
	}

	if errors.Is(err, io.EOF) {
		return commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "body",
			Constraints: []string{"body must not be empty"},
		}})
	}

	return commonapi.NewInternalServerError()
}

func validationDetails(errs validator.ValidationErrors) []commonapi.ErrorDetail {
	details := make([]commonapi.ErrorDetail, 0, len(errs))
	indexByField := make(map[string]int, len(errs))

	for _, validationErr := range errs {
		field := validationErr.Field()
		if field == "" {
			field = "body"
		}

		message := validationConstraint(field, validationErr)

		if index, exists := indexByField[field]; exists {
			details[index].Constraints = append(details[index].Constraints, message)
			continue
		}

		indexByField[field] = len(details)
		details = append(details, commonapi.ErrorDetail{
			Field:       field,
			Constraints: []string{message},
		})
	}

	return details
}

func validationConstraint(field string, validationErr validator.FieldError) string {
	switch validationErr.Tag() {
	case "required":
		return fmt.Sprintf("%s must not be empty", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, validationErr.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters long", field, validationErr.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of %s", field, strings.ReplaceAll(validationErr.Param(), " ", ", "))
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, validationErr.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, validationErr.Param())
	default:
		return fmt.Sprintf("%s failed validation on %s", field, validationErr.Tag())
	}
}

func humanizeType(valueType reflect.Type) string {
	if valueType == nil {
		return "the expected type"
	}

	switch valueType.Kind() {
	case reflect.Bool:
		return "a boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "an integer"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "a positive integer"
	case reflect.Float32, reflect.Float64:
		return "a number"
	case reflect.String:
		return "a string"
	case reflect.Array, reflect.Slice:
		return "an array"
	case reflect.Map, reflect.Struct:
		return "an object"
	default:
		return "the expected type"
	}
}
