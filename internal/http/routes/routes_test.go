package routes_test

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
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

func TestVersionedHealthRouteReturnsOKStatus(t *testing.T) {
	t.Parallel()

	router := routes.New(config.Config{AppEnv: "test"}, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for versioned health route, got %d", recorder.Code)
	}

	var response commonapi.StatusResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("expected status to be ok, got %q", response.Status)
	}
}

func TestRootRouteReturnsDeprecationSignal(t *testing.T) {
	t.Parallel()

	router := routes.New(config.Config{AppEnv: "test"}, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for deprecated root route, got %d", recorder.Code)
	}

	if got := recorder.Header().Get("Deprecation"); got != "true" {
		t.Fatalf("expected Deprecation header to be true, got %q", got)
	}

	var response commonapi.MessageResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	want := "The unversioned root route is deprecated. Migrate to /v1/health."
	if response.Message != want {
		t.Fatalf("expected deprecation message %q, got %q", want, response.Message)
	}
}

func TestUnversionedHealthRouteIsNotRegistered(t *testing.T) {
	t.Parallel()

	router := routes.New(config.Config{AppEnv: "test"}, slog.Default())
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unversioned health route, got %d", recorder.Code)
	}
}
