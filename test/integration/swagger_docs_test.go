package integration_test

import (
	"net/http"
	"strings"
	"testing"

	"go-gin-ecommerce/test/integration/testutil"
)

func TestSwaggerUIAndSpecAreServed(t *testing.T) {
	t.Parallel()

	router := testutil.NewRouter()

	uiRecorder := performRequestWithHeaders(t, router, http.MethodGet, "/swagger/index.html", "", nil)
	if uiRecorder.Code != http.StatusOK {
		t.Fatalf("expected swagger UI 200, got %d", uiRecorder.Code)
	}
	if !strings.Contains(uiRecorder.Body.String(), "Swagger UI") {
		t.Fatalf("expected swagger UI page, got body %q", uiRecorder.Body.String())
	}

	specRecorder := performRequestWithHeaders(t, router, http.MethodGet, "/swagger/doc.json", "", nil)
	if specRecorder.Code != http.StatusOK {
		t.Fatalf("expected swagger spec 200, got %d", specRecorder.Code)
	}

	body := specRecorder.Body.String()
	for _, expected := range []string{`"/products"`, `"/orders"`, `"/orders/{id}/cancel"`} {
		if !strings.Contains(body, expected) {
			t.Fatalf("expected swagger spec to contain %s", expected)
		}
	}
	if !strings.Contains(body, `"sku"`) {
		t.Fatal("expected swagger spec to mention deprecated sku alias")
	}
}
