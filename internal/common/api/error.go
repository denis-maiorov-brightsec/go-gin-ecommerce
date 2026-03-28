package api

import (
	"net/http"
	"time"
)

const timestampLayout = "2006-01-02T15:04:05.000Z07:00"

type ErrorDetail struct {
	Field       string   `json:"field,omitempty" example:"stockKeepingUnit"`
	Constraints []string `json:"constraints,omitempty" example:"stockKeepingUnit must not be empty"`
}

type ErrorPayload struct {
	Code    string        `json:"code" example:"VALIDATION_ERROR"`
	Message string        `json:"message" example:"Request validation failed"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorResponse struct {
	Timestamp string       `json:"timestamp" example:"2026-03-28T12:34:56.000Z"`
	Path      string       `json:"path" example:"/v1/products"`
	Error     ErrorPayload `json:"error"`
}

type Error struct {
	Status  int
	Code    string
	Message string
	Details []ErrorDetail
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}

func NewValidationError(details []ErrorDetail) *Error {
	return &Error{
		Status:  http.StatusBadRequest,
		Code:    "VALIDATION_ERROR",
		Message: "Request validation failed",
		Details: details,
	}
}

func NewNotFoundError() *Error {
	return &Error{
		Status:  http.StatusNotFound,
		Code:    "NOT_FOUND",
		Message: "Resource not found",
	}
}

func NewConflictError(message string, details []ErrorDetail) *Error {
	if message == "" {
		message = "Resource conflict"
	}

	return &Error{
		Status:  http.StatusConflict,
		Code:    "CONFLICT",
		Message: message,
		Details: details,
	}
}

func NewUnauthorizedError() *Error {
	return &Error{
		Status:  http.StatusUnauthorized,
		Code:    "UNAUTHORIZED",
		Message: "Authentication required",
	}
}

func NewForbiddenError() *Error {
	return &Error{
		Status:  http.StatusForbidden,
		Code:    "FORBIDDEN",
		Message: "Forbidden",
	}
}

func NewTooManyRequestsError() *Error {
	return &Error{
		Status:  http.StatusTooManyRequests,
		Code:    "TOO_MANY_REQUESTS",
		Message: "Rate limit exceeded",
	}
}

func NewInternalServerError() *Error {
	return &Error{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: "Internal server error",
	}
}

func NewErrorResponse(path string, apiErr *Error) ErrorResponse {
	if apiErr == nil {
		apiErr = NewInternalServerError()
	}

	return ErrorResponse{
		Timestamp: time.Now().UTC().Format(timestampLayout),
		Path:      path,
		Error: ErrorPayload{
			Code:    apiErr.Code,
			Message: apiErr.Message,
			Details: apiErr.Details,
		},
	}
}
