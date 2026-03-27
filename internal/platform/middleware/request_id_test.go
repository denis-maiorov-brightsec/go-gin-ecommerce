package middleware_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-ecommerce/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

func TestRequestIDPropagatesIncomingHeaderToHandlerAndResponse(t *testing.T) {
	t.Parallel()

	router := newRequestIDTestRouter(slog.Default())
	router.GET("/v1/test", func(c *gin.Context) {
		requestID, ok := middleware.RequestIDFromContext(c.Request.Context())
		if !ok {
			t.Fatal("expected request id in request context")
		}

		c.JSON(http.StatusOK, gin.H{
			"requestId": requestID,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	req.Header.Set(middleware.RequestIDHeader, "req-from-client")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	if got := recorder.Header().Get(middleware.RequestIDHeader); got != "req-from-client" {
		t.Fatalf("expected response request id to echo incoming header, got %q", got)
	}

	var response struct {
		RequestID string `json:"requestId"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.RequestID != "req-from-client" {
		t.Fatalf("expected handler to receive propagated request id, got %q", response.RequestID)
	}
}

func TestRequestIDGeneratesIDWhenHeaderMissing(t *testing.T) {
	t.Parallel()

	router := newRequestIDTestRouter(slog.Default())
	router.GET("/v1/test", func(c *gin.Context) {
		requestID, ok := middleware.RequestIDFromContext(c.Request.Context())
		if !ok {
			t.Fatal("expected generated request id in request context")
		}

		c.JSON(http.StatusOK, gin.H{
			"requestId": requestID,
		})
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	generatedID := recorder.Header().Get(middleware.RequestIDHeader)
	if generatedID == "" {
		t.Fatal("expected generated request id header")
	}

	var response struct {
		RequestID string `json:"requestId"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.RequestID != generatedID {
		t.Fatalf("expected generated request id %q in handler response, got %q", generatedID, response.RequestID)
	}
}

func TestRequestLoggerEmitsStructuredLogWithRequestID(t *testing.T) {
	t.Parallel()

	var buffer bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buffer, nil))
	router := newRequestIDTestRouter(logger)
	router.GET("/v1/test", func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/test", nil)
	req.Header.Set(middleware.RequestIDHeader, "req-log-123")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", recorder.Code)
	}

	var entry map[string]any
	if err := json.Unmarshal(buffer.Bytes(), &entry); err != nil {
		t.Fatalf("failed to decode log entry: %v", err)
	}

	if entry["msg"] != "request completed" {
		t.Fatalf("expected request completed log message, got %#v", entry["msg"])
	}

	if entry["request_id"] != "req-log-123" {
		t.Fatalf("expected logged request id req-log-123, got %#v", entry["request_id"])
	}

	if entry["method"] != http.MethodGet {
		t.Fatalf("expected method GET, got %#v", entry["method"])
	}

	if entry["path"] != "/v1/test" {
		t.Fatalf("expected path /v1/test, got %#v", entry["path"])
	}

	if status, ok := entry["status"].(float64); !ok || int(status) != http.StatusNoContent {
		t.Fatalf("expected status 204, got %#v", entry["status"])
	}

	if _, ok := entry["latency_ms"].(float64); !ok {
		t.Fatalf("expected latency_ms field, got %#v", entry["latency_ms"])
	}
}

func newRequestIDTestRouter(logger *slog.Logger) *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(middleware.RequestID())
	router.Use(middleware.RequestLogger(logger))
	router.Use(middleware.Recovery(logger))
	router.Use(middleware.ErrorHandler(logger))

	return router
}

func TestRequestIDFromContextReturnsFalseWhenMissing(t *testing.T) {
	t.Parallel()

	if requestID, ok := middleware.RequestIDFromContext(context.Background()); ok || requestID != "" {
		t.Fatalf("expected missing request id, got %q, %t", requestID, ok)
	}
}
