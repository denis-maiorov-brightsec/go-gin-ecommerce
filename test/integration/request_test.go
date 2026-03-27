package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-gin-ecommerce/internal/platform/middleware"
	"go-gin-ecommerce/test/integration/testutil"
)

func performRequestWithHeaders(t *testing.T, router http.Handler, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	var requestBody *strings.Reader
	if body == "" {
		requestBody = strings.NewReader("")
	} else {
		requestBody = strings.NewReader(body)
	}

	req := httptest.NewRequest(method, path, requestBody)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func performJSONRequestWithHeaders(t *testing.T, router http.Handler, method string, path string, body string, headers map[string]string) *httptest.ResponseRecorder {
	t.Helper()

	return performRequestWithHeaders(t, router, method, path, body, headers)
}

func TestHealthRouteEchoesIncomingRequestIDHeader(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()

	recorder := performRequestWithHeaders(t, router, http.MethodGet, "/v1/health", "", map[string]string{
		middleware.RequestIDHeader: "req-health-echo",
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	if got := recorder.Header().Get(middleware.RequestIDHeader); got != "req-health-echo" {
		t.Fatalf("expected echoed request id header, got %q", got)
	}
}

func TestHealthRouteGeneratesRequestIDHeaderWhenMissing(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()

	recorder := performRequestWithHeaders(t, router, http.MethodGet, "/v1/health", "", nil)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}

	requestID := recorder.Header().Get(middleware.RequestIDHeader)
	if requestID == "" {
		t.Fatal("expected generated request id header")
	}

	var response map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode health response: %v", err)
	}

	if response["status"] != "ok" {
		t.Fatalf("expected ok status, got %#v", response["status"])
	}
}
