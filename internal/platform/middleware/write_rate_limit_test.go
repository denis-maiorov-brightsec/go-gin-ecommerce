package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gin-ecommerce/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

func TestWriteRateLimiterRejectsRequestsOverLimit(t *testing.T) {
	t.Parallel()

	currentTime := time.Date(2026, time.March, 27, 12, 0, 0, 0, time.UTC)
	router := newTestRouter()
	router.POST("/v1/test", middleware.NewWriteRateLimiter(middleware.WriteRateLimiterConfig{
		Limit:  2,
		Window: time.Minute,
		Now: func() time.Time {
			return currentTime
		},
	}), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	for attempt := 1; attempt <= 2; attempt++ {
		req := httptest.NewRequest(http.MethodPost, "/v1/test", nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, req)

		if recorder.Code != http.StatusNoContent {
			t.Fatalf("expected request %d to succeed, got %d with body %s", attempt, recorder.Code, recorder.Body.String())
		}
	}

	req := httptest.NewRequest(http.MethodPost, "/v1/test", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 once limit is exceeded, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeErrorResponse(t, recorder)
	assertErrorPayload(t, response, "TOO_MANY_REQUESTS", "Rate limit exceeded")
}

func TestWriteRateLimiterResetsAfterWindow(t *testing.T) {
	t.Parallel()

	currentTime := time.Date(2026, time.March, 27, 12, 0, 0, 0, time.UTC)
	router := newTestRouter()
	router.POST("/v1/test", middleware.NewWriteRateLimiter(middleware.WriteRateLimiterConfig{
		Limit:  1,
		Window: time.Minute,
		Now: func() time.Time {
			return currentTime
		},
	}), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	firstReq := httptest.NewRequest(http.MethodPost, "/v1/test", nil)
	firstRecorder := httptest.NewRecorder()
	router.ServeHTTP(firstRecorder, firstReq)
	if firstRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected first request to succeed, got %d", firstRecorder.Code)
	}

	currentTime = currentTime.Add(time.Minute)

	secondReq := httptest.NewRequest(http.MethodPost, "/v1/test", nil)
	secondRecorder := httptest.NewRecorder()
	router.ServeHTTP(secondRecorder, secondReq)
	if secondRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected request after window reset to succeed, got %d with body %s", secondRecorder.Code, secondRecorder.Body.String())
	}
}
