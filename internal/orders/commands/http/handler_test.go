package http

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/orders/model"
	"go-gin-ecommerce/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

func TestCancelRejectsInvalidOrderID(t *testing.T) {
	t.Helper()

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(middleware.ErrorHandler(slog.New(slog.NewTextHandler(io.Discard, nil))))

	handler := NewHandler(&stubService{})
	handler.RegisterRoutes(router.Group("/orders"))

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/orders/nope/cancel", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid order id, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %#v", response.Error)
	}
}

type stubService struct{}

func (s *stubService) Cancel(ctx context.Context, id uint) (model.Order, error) {
	return model.Order{}, nil
}
