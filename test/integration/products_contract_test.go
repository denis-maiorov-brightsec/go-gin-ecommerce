package integration_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"

	"go-gin-ecommerce/test/integration/testutil"
)

func TestProductsContractCreateGetUpdateDeleteFlow(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	createRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
		"name": "Contract Lamp",
		"sku": "CONTRACT-001",
		"price": 49.99,
		"status": "active"
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating product, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	created := decodeJSONMap(t, createRecorder.Body.Bytes())
	assertJSONKeys(t, created, "created product", "createdAt", "id", "name", "price", "status", "stockKeepingUnit", "updatedAt")
	assertNoJSONKey(t, created, "sku")
	assertJSONNumberValue(t, created, "id", 1)
	assertJSONStringValue(t, created, "name", "Contract Lamp")
	assertJSONStringValue(t, created, "stockKeepingUnit", "CONTRACT-001")
	assertJSONNumberValue(t, created, "price", 49.99)
	assertJSONStringValue(t, created, "status", "active")
	assertJSONStringFieldPresent(t, created, "createdAt")
	assertJSONStringFieldPresent(t, created, "updatedAt")

	getRecorder := performRequest(t, router, http.MethodGet, "/v1/products/1", "")
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when fetching product, got %d with body %s", getRecorder.Code, getRecorder.Body.String())
	}

	got := decodeJSONMap(t, getRecorder.Body.Bytes())
	assertJSONKeys(t, got, "fetched product", "createdAt", "id", "name", "price", "status", "stockKeepingUnit", "updatedAt")
	assertNoJSONKey(t, got, "sku")
	assertJSONNumberValue(t, got, "id", 1)
	assertJSONStringValue(t, got, "name", "Contract Lamp")
	assertJSONStringValue(t, got, "stockKeepingUnit", "CONTRACT-001")

	updateRecorder := performJSONRequest(t, router, http.MethodPatch, "/v1/products/1", `{
		"name": "Contract Lamp v2",
		"stockKeepingUnit": "CONTRACT-002",
		"price": 59.99,
		"status": "draft",
		"categoryId": 12
	}`)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when updating product, got %d with body %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	updated := decodeJSONMap(t, updateRecorder.Body.Bytes())
	assertJSONKeys(t, updated, "updated product", "categoryId", "createdAt", "id", "name", "price", "status", "stockKeepingUnit", "updatedAt")
	assertNoJSONKey(t, updated, "sku")
	assertJSONStringValue(t, updated, "name", "Contract Lamp v2")
	assertJSONStringValue(t, updated, "stockKeepingUnit", "CONTRACT-002")
	assertJSONNumberValue(t, updated, "price", 59.99)
	assertJSONStringValue(t, updated, "status", "draft")
	assertJSONNumberValue(t, updated, "categoryId", 12)

	deleteRecorder := performRequest(t, router, http.MethodDelete, "/v1/products/1", "")
	if deleteRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 when deleting product, got %d with body %q", deleteRecorder.Code, deleteRecorder.Body.String())
	}
	if deleteRecorder.Body.Len() != 0 {
		t.Fatalf("expected empty body for delete response, got %q", deleteRecorder.Body.String())
	}

	missingRecorder := performRequest(t, router, http.MethodGet, "/v1/products/1", "")
	if missingRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deleting product, got %d with body %s", missingRecorder.Code, missingRecorder.Body.String())
	}

	missing := decodeJSONMap(t, missingRecorder.Body.Bytes())
	assertErrorEnvelope(t, missing, "/v1/products/1", "NOT_FOUND", "Resource not found")
}

func TestProductsContractListResponseShape(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	for i, body := range []string{
		`{"name":"Contract Product 1","stockKeepingUnit":"LIST-001","price":10,"status":"active"}`,
		`{"name":"Contract Product 2","stockKeepingUnit":"LIST-002","price":20,"status":"draft"}`,
	} {
		recorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", body)
		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected 201 when seeding product %d, got %d with body %s", i+1, recorder.Code, recorder.Body.String())
		}
	}

	recorder := performRequest(t, router, http.MethodGet, "/v1/products?page=1&limit=1", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when listing products, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeJSONMap(t, recorder.Body.Bytes())
	assertJSONKeys(t, response, "products list", "items", "limit", "page", "total", "totalPages")
	assertJSONNumberValue(t, response, "page", 1)
	assertJSONNumberValue(t, response, "limit", 1)
	assertJSONNumberValue(t, response, "total", 2)
	assertJSONNumberValue(t, response, "totalPages", 2)

	items := getJSONArray(t, response, "items")
	if len(items) != 1 {
		t.Fatalf("expected 1 listed item, got %d", len(items))
	}

	first := items[0]
	assertJSONKeys(t, first, "listed product", "createdAt", "id", "name", "price", "status", "stockKeepingUnit", "updatedAt")
	assertNoJSONKey(t, first, "sku")
	assertJSONStringValue(t, first, "name", "Contract Product 1")
	assertJSONStringValue(t, first, "stockKeepingUnit", "LIST-001")
}

func TestProductsContractSearchResponseShape(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	for i, body := range []string{
		`{"name":"Desk Search Result","stockKeepingUnit":"SEARCH-001","price":19.99,"status":"active"}`,
		`{"name":"Desk Search Result 2","stockKeepingUnit":"SEARCH-002","price":29.99,"status":"draft"}`,
		`{"name":"Other Result","stockKeepingUnit":"OTHER-003","price":39.99,"status":"active"}`,
	} {
		recorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", body)
		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected 201 when seeding search product %d, got %d with body %s", i+1, recorder.Code, recorder.Body.String())
		}
	}

	recorder := performRequest(t, router, http.MethodGet, "/v1/search/products?q=desk", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when searching products, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeJSONMap(t, recorder.Body.Bytes())
	assertJSONKeys(t, response, "search response", "items", "limit", "page", "total", "totalPages")
	assertJSONNumberValue(t, response, "page", 1)
	assertJSONNumberValue(t, response, "limit", 20)
	assertJSONNumberValue(t, response, "total", 2)
	assertJSONNumberValue(t, response, "totalPages", 1)

	items := getJSONArray(t, response, "items")
	if len(items) != 2 {
		t.Fatalf("expected 2 search items, got %d", len(items))
	}

	for i, item := range items {
		assertJSONKeys(t, item, "search item", "createdAt", "id", "name", "price", "status", "stockKeepingUnit", "updatedAt")
		assertNoJSONKey(t, item, "sku")
		assertJSONStringFieldPresent(t, item, "stockKeepingUnit")
		if i == 0 {
			assertJSONStringValue(t, item, "name", "Desk Search Result")
		}
	}
}

func TestProductsContractValidationErrorEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	createRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
		"name": "Validation Seed",
		"stockKeepingUnit": "VALID-001",
		"price": 10,
		"status": "active"
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when seeding product, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	testCases := []struct {
		name            string
		method          string
		path            string
		body            string
		expectedField   string
		expectedMessage string
	}{
		{
			name:            "list invalid pagination",
			method:          http.MethodGet,
			path:            "/v1/products?page=0",
			expectedField:   "page",
			expectedMessage: "Request validation failed",
		},
		{
			name:            "create conflicting canonical and alias sku",
			method:          http.MethodPost,
			path:            "/v1/products",
			body:            `{"name":"Conflict","stockKeepingUnit":"VALID-002","sku":"DIFFERENT-002","price":11,"status":"active"}`,
			expectedField:   "stockKeepingUnit",
			expectedMessage: "Request validation failed",
		},
		{
			name:            "patch invalid price",
			method:          http.MethodPatch,
			path:            "/v1/products/1",
			body:            `{"price":-1}`,
			expectedField:   "price",
			expectedMessage: "Request validation failed",
		},
		{
			name:            "search missing query",
			method:          http.MethodGet,
			path:            "/v1/search/products",
			expectedField:   "q",
			expectedMessage: "Request validation failed",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			recorder := performJSONRequest(t, router, tt.method, tt.path, tt.body)
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected 400, got %d with body %s", recorder.Code, recorder.Body.String())
			}

			response := decodeJSONMap(t, recorder.Body.Bytes())
			assertErrorEnvelope(t, response, trimQuery(tt.path), "VALIDATION_ERROR", tt.expectedMessage)

			errorPayload := getJSONObject(t, response, "error")
			details := getJSONArray(t, errorPayload, "details")
			if len(details) != 1 {
				t.Fatalf("expected 1 validation detail, got %d", len(details))
			}

			assertJSONKeys(t, details[0], "error detail", "constraints", "field")
			assertJSONStringValue(t, details[0], "field", tt.expectedField)
		})
	}
}

func TestProductsContractNotFoundErrorEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	testCases := []struct {
		name   string
		method string
		path   string
		body   string
	}{
		{name: "get missing product", method: http.MethodGet, path: "/v1/products/999"},
		{name: "patch missing product", method: http.MethodPatch, path: "/v1/products/999", body: `{"status":"draft"}`},
		{name: "delete missing product", method: http.MethodDelete, path: "/v1/products/999"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			recorder := performJSONRequest(t, router, tt.method, tt.path, tt.body)
			if recorder.Code != http.StatusNotFound {
				t.Fatalf("expected 404, got %d with body %s", recorder.Code, recorder.Body.String())
			}

			response := decodeJSONMap(t, recorder.Body.Bytes())
			assertErrorEnvelope(t, response, tt.path, "NOT_FOUND", "Resource not found")
		})
	}
}

func decodeJSONMap(t *testing.T, body []byte) map[string]any {
	t.Helper()

	var decoded map[string]any
	if err := json.Unmarshal(body, &decoded); err != nil {
		t.Fatalf("failed to decode json payload: %v", err)
	}

	return decoded
}

func assertErrorEnvelope(t *testing.T, payload map[string]any, path string, code string, message string) {
	t.Helper()

	assertJSONKeys(t, payload, "error response", "error", "path", "timestamp")
	assertJSONStringFieldPresent(t, payload, "timestamp")
	assertJSONStringValue(t, payload, "path", path)

	errorPayload := getJSONObject(t, payload, "error")
	assertJSONStringValue(t, errorPayload, "code", code)
	assertJSONStringValue(t, errorPayload, "message", message)

	switch code {
	case "VALIDATION_ERROR":
		assertJSONKeys(t, errorPayload, "error payload", "code", "details", "message")
	default:
		assertJSONKeys(t, errorPayload, "error payload", "code", "message")
	}
}

func assertJSONKeys(t *testing.T, payload map[string]any, label string, expectedKeys ...string) {
	t.Helper()

	actualKeys := make([]string, 0, len(payload))
	for key := range payload {
		actualKeys = append(actualKeys, key)
	}

	if !sameStringSet(actualKeys, expectedKeys) {
		t.Fatalf("unexpected keys for %s: got %v want %v", label, actualKeys, expectedKeys)
	}
}

func sameStringSet(actual []string, expected []string) bool {
	actualSet := make(map[string]struct{}, len(actual))
	for _, value := range actual {
		actualSet[value] = struct{}{}
	}

	expectedSet := make(map[string]struct{}, len(expected))
	for _, value := range expected {
		expectedSet[value] = struct{}{}
	}

	return reflect.DeepEqual(actualSet, expectedSet)
}

func assertNoJSONKey(t *testing.T, payload map[string]any, key string) {
	t.Helper()

	if _, exists := payload[key]; exists {
		t.Fatalf("expected key %q to be absent in %#v", key, payload)
	}
}

func assertJSONStringValue(t *testing.T, payload map[string]any, key string, expected string) {
	t.Helper()

	value, ok := payload[key].(string)
	if !ok {
		t.Fatalf("expected %q to be a string, got %#v", key, payload[key])
	}
	if value != expected {
		t.Fatalf("unexpected value for %q: got %q want %q", key, value, expected)
	}
}

func assertJSONStringFieldPresent(t *testing.T, payload map[string]any, key string) {
	t.Helper()

	value, ok := payload[key].(string)
	if !ok || value == "" {
		t.Fatalf("expected non-empty string for %q, got %#v", key, payload[key])
	}
}

func assertJSONNumberValue(t *testing.T, payload map[string]any, key string, expected float64) {
	t.Helper()

	value, ok := payload[key].(float64)
	if !ok {
		t.Fatalf("expected %q to be a number, got %#v", key, payload[key])
	}
	if value != expected {
		t.Fatalf("unexpected numeric value for %q: got %v want %v", key, value, expected)
	}
}

func getJSONObject(t *testing.T, payload map[string]any, key string) map[string]any {
	t.Helper()

	value, ok := payload[key].(map[string]any)
	if !ok {
		t.Fatalf("expected %q to be an object, got %#v", key, payload[key])
	}

	return value
}

func getJSONArray(t *testing.T, payload map[string]any, key string) []map[string]any {
	t.Helper()

	raw, ok := payload[key].([]any)
	if !ok {
		t.Fatalf("expected %q to be an array, got %#v", key, payload[key])
	}

	items := make([]map[string]any, 0, len(raw))
	for _, entry := range raw {
		object, ok := entry.(map[string]any)
		if !ok {
			t.Fatalf("expected array entry for %q to be an object, got %#v", key, entry)
		}
		items = append(items, object)
	}

	return items
}

func trimQuery(path string) string {
	for i, char := range path {
		if char == '?' {
			return path[:i]
		}
	}

	return path
}
