package api

import (
	"net/http"
	"time"
)

const timestampLayout = "2006-01-02T15:04:05.000Z07:00"

type ErrorDetail struct {
	Field       string   `json:"field,omitempty"`
	Constraints []string `json:"constraints,omitempty"`
}

type ErrorPayload struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type ErrorResponse struct {
	Timestamp string       `json:"timestamp"`
	Path      string       `json:"path"`
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
