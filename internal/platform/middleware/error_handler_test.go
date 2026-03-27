package middleware_test

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

const timestampLayout = "2006-01-02T15:04:05.000Z07:00"

func TestErrorHandlerFormatsValidationErrors(t *testing.T) {
	t.Parallel()

	router := newTestRouter()
	router.POST("/v1/test", func(c *gin.Context) {
		var payload struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/test", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for validation failure, got %d", recorder.Code)
	}

	response := decodeErrorResponse(t, recorder)

	if response.Path != "/v1/test" {
		t.Fatalf("expected path to be /v1/test, got %q", response.Path)
	}

	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}

	if response.Error.Message != "Request validation failed" {
		t.Fatalf("expected validation error message, got %q", response.Error.Message)
	}

	if len(response.Error.Details) != 1 {
		t.Fatalf("expected a single validation detail, got %d", len(response.Error.Details))
	}

	detail := response.Error.Details[0]
	if detail.Field != "name" {
		t.Fatalf("expected validation detail field to be name, got %q", detail.Field)
	}

	if len(detail.Constraints) != 1 || detail.Constraints[0] != "name must not be empty" {
		t.Fatalf("expected required validation constraint, got %#v", detail.Constraints)
	}
}

func TestErrorHandlerSanitizesUnhandledErrors(t *testing.T) {
	t.Parallel()

	router := newTestRouter()
	router.GET("/v1/test", func(c *gin.Context) {
		_ = c.Error(errors.New("database connection refused"))
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for unhandled error, got %d", recorder.Code)
	}

	response := decodeErrorResponse(t, recorder)

	if response.Error.Code != "INTERNAL_SERVER_ERROR" {
		t.Fatalf("expected internal server error code, got %q", response.Error.Code)
	}

	if response.Error.Message != "Internal server error" {
		t.Fatalf("expected sanitized internal server error message, got %q", response.Error.Message)
	}

	if strings.Contains(recorder.Body.String(), "database connection refused") {
		t.Fatalf("expected response body to omit original error, got %s", recorder.Body.String())
	}
}

func TestRecoveryFormatsPanicErrors(t *testing.T) {
	t.Parallel()

	router := newTestRouter()
	router.GET("/v1/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/panic", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for panic recovery, got %d", recorder.Code)
	}

	response := decodeErrorResponse(t, recorder)

	if response.Path != "/v1/panic" {
		t.Fatalf("expected path to be /v1/panic, got %q", response.Path)
	}

	if response.Error.Code != "INTERNAL_SERVER_ERROR" {
		t.Fatalf("expected internal server error code, got %q", response.Error.Code)
	}
}

func newTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	middleware.SetupValidation()

	router := gin.New()
	router.Use(middleware.Recovery(slog.Default()))
	router.Use(middleware.ErrorHandler(slog.Default()))

	return router
}

func decodeErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) commonapi.ErrorResponse {
	t.Helper()

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	if _, err := time.Parse(timestampLayout, response.Timestamp); err != nil {
		t.Fatalf("expected timestamp in RFC3339-with-millis format, got %q: %v", response.Timestamp, err)
	}

	return response
}
