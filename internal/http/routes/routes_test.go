package routes_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-ecommerce/internal/http/routes"
	"go-gin-ecommerce/internal/platform/config"
)

func TestNewRouterBuildsAndReturnsNotFoundForUnknownRoute(t *testing.T) {
	t.Parallel()

	router := routes.New(config.Config{AppEnv: "test"}, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/does-not-exist", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown route, got %d", recorder.Code)
	}
}
