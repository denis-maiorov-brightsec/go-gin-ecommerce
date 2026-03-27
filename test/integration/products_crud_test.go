package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/products/dto"
	"go-gin-ecommerce/test/integration/testutil"
)

func TestProductsCRUD(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	createRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
		"name": "Desk Lamp",
		"sku": "LAMP-001",
		"price": 49.99,
		"status": "active"
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating product, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	created := decodeProductResponse(t, createRecorder)
	if created.ID == 0 {
		t.Fatalf("expected created product to have an id, got %#v", created)
	}
	if created.Name != "Desk Lamp" || created.SKU != "LAMP-001" || created.Price != 49.99 || created.Status != "active" {
		t.Fatalf("unexpected create response: %#v", created)
	}

	listRecorder := performRequest(t, router, http.MethodGet, "/v1/products", "")
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when listing products, got %d", listRecorder.Code)
	}

	var listed []dto.ProductResponse
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &listed); err != nil {
		t.Fatalf("failed to decode product list: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != created.ID {
		t.Fatalf("unexpected listed products: %#v", listed)
	}

	getRecorder := performRequest(t, router, http.MethodGet, "/v1/products/1", "")
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when fetching product, got %d", getRecorder.Code)
	}

	updateRecorder := performJSONRequest(t, router, http.MethodPatch, "/v1/products/1", `{
		"price": 59.99,
		"status": "draft"
	}`)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when updating product, got %d with body %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	updated := decodeProductResponse(t, updateRecorder)
	if updated.Price != 59.99 || updated.Status != "draft" || updated.Name != created.Name || updated.SKU != created.SKU {
		t.Fatalf("unexpected updated product: %#v", updated)
	}

	deleteRecorder := performRequest(t, router, http.MethodDelete, "/v1/products/1", "")
	if deleteRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 when deleting product, got %d", deleteRecorder.Code)
	}

	missingRecorder := performRequest(t, router, http.MethodGet, "/v1/products/1", "")
	if missingRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deleting product, got %d", missingRecorder.Code)
	}
}

func TestCreateProductValidationErrorUsesGlobalEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
		"name": "",
		"sku": "BAD-001",
		"price": 0,
		"status": ""
	}`)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid create request, got %d", recorder.Code)
	}

	response := decodeProductErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
	if response.Path != "/v1/products" {
		t.Fatalf("expected response path /v1/products, got %q", response.Path)
	}
}

func TestGetMissingProductReturnsNotFoundEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/products/999", "")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing product, got %d", recorder.Code)
	}

	response := decodeProductErrorResponse(t, recorder)
	if response.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND error code, got %q", response.Error.Code)
	}
}

func TestDeleteMissingProductReturnsNotFoundEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequest(t, router, http.MethodDelete, "/v1/products/999", "")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleting missing product, got %d", recorder.Code)
	}

	response := decodeProductErrorResponse(t, recorder)
	if response.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND error code, got %q", response.Error.Code)
	}
}

func TestPatchProductValidationErrorUsesGlobalEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	createRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
		"name": "Notebook",
		"sku": "NOTE-001",
		"price": 9.99,
		"status": "active"
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating product, got %d", createRecorder.Code)
	}

	recorder := performJSONRequest(t, router, http.MethodPatch, "/v1/products/1", `{
		"price": -1
	}`)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid patch request, got %d", recorder.Code)
	}

	response := decodeProductErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
}

func TestPatchProductCanClearCategoryID(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	createRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
		"name": "Monitor Stand",
		"sku": "STAND-001",
		"price": 39.99,
		"status": "active",
		"categoryId": 7
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating product, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	created := decodeProductResponse(t, createRecorder)
	if created.CategoryID == nil || *created.CategoryID != 7 {
		t.Fatalf("expected created product categoryId to be 7, got %#v", created.CategoryID)
	}

	updateRecorder := performJSONRequest(t, router, http.MethodPatch, "/v1/products/1", `{
		"categoryId": null
	}`)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when clearing categoryId, got %d with body %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	updated := decodeProductResponse(t, updateRecorder)
	if updated.CategoryID != nil {
		t.Fatalf("expected categoryId to be cleared, got %#v", updated.CategoryID)
	}
}

func performRequest(t *testing.T, router http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
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

	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)

	return recorder
}

func performJSONRequest(t *testing.T, router http.Handler, method string, path string, body string) *httptest.ResponseRecorder {
	t.Helper()
	return performRequest(t, router, method, path, body)
}

func decodeProductResponse(t *testing.T, recorder *httptest.ResponseRecorder) dto.ProductResponse {
	t.Helper()

	var response dto.ProductResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode product response: %v", err)
	}

	return response
}

func decodeProductErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) commonapi.ErrorResponse {
	t.Helper()

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	return response
}
