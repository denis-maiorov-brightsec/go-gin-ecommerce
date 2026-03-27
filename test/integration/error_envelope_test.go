package integration_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/test/integration/testutil"

	"github.com/gin-gonic/gin"
)

func TestValidationErrorsUseGlobalEnvelope(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()
	router.POST("/v1/test-validation", func(c *gin.Context) {
		var payload struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.ShouldBindJSON(&payload); err != nil {
			_ = c.Error(err)
			return
		}

		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/test-validation", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payload, got %d", recorder.Code)
	}

	response := decodeIntegrationErrorResponse(t, recorder)

	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}

	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "name" {
		t.Fatalf("expected validation details for name, got %#v", response.Error.Details)
	}
}

func TestBindJSONValidationErrorsUseGlobalEnvelope(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()
	router.POST("/v1/test-validation", func(c *gin.Context) {
		var payload struct {
			Name string `json:"name" binding:"required"`
		}

		if err := c.BindJSON(&payload); err != nil {
			return
		}

		c.Status(http.StatusCreated)
	})

	req := httptest.NewRequest(http.MethodPost, "/v1/test-validation", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid payload, got %d", recorder.Code)
	}

	response := decodeIntegrationErrorResponse(t, recorder)

	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}

	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "name" {
		t.Fatalf("expected validation details for name, got %#v", response.Error.Details)
	}
}

func TestInternalErrorsUseSanitizedGlobalEnvelope(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()
	router.GET("/v1/test-error", func(c *gin.Context) {
		_ = c.Error(errors.New("raw database failure"))
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test-error", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for runtime error, got %d", recorder.Code)
	}

	response := decodeIntegrationErrorResponse(t, recorder)

	if response.Error.Message != "Internal server error" {
		t.Fatalf("expected sanitized message, got %q", response.Error.Message)
	}

	if strings.Contains(recorder.Body.String(), "raw database failure") {
		t.Fatalf("expected response body to omit raw runtime error, got %s", recorder.Body.String())
	}
}

func decodeIntegrationErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) commonapi.ErrorResponse {
	t.Helper()

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	return response
}
