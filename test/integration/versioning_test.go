package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/test/integration/testutil"
)

func TestVersionedHealthEndpoint(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/v1/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for /v1/health, got %d", recorder.Code)
	}

	var response commonapi.StatusResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "ok" {
		t.Fatalf("expected status to be ok, got %q", response.Status)
	}
}

func TestDeprecatedRootEndpoint(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for /, got %d", recorder.Code)
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

func TestUnversionedHealthEndpointIsNotAvailable(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for /health, got %d", recorder.Code)
	}

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Path != "/health" {
		t.Fatalf("expected response path to be /health, got %q", response.Path)
	}

	if response.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected error code NOT_FOUND, got %q", response.Error.Code)
	}
}
