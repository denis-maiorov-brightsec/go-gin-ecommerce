package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/products/dto"
	"go-gin-ecommerce/test/integration/testutil"
)

func TestSearchProductsMatchesNameAndSKUCaseInsensitivelyWithDeterministicOrdering(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	for i, body := range []string{
		`{"name":"Omega Lamp","sku":"DESK-LIGHT-003","price":59.99,"status":"active"}`,
		`{"name":"alpha desk","sku":"WORK-001","price":39.99,"status":"active"}`,
		`{"name":"Beta Desk Shelf","sku":"SHELF-002","price":79.99,"status":"draft"}`,
		`{"name":"Chair Mat","sku":"FLOOR-004","price":29.99,"status":"active"}`,
	} {
		recorder := performJSONRequest(t, router, http.MethodPost, "/v1/products", body)
		if recorder.Code != http.StatusCreated {
			t.Fatalf("expected 201 when creating product %d, got %d with body %s", i+1, recorder.Code, recorder.Body.String())
		}
	}

	recorder := performRequest(t, router, http.MethodGet, "/v1/search/products?q=desk", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when searching products, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response dto.ProductListResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode search response: %v", err)
	}

	if response.Page != 1 || response.Limit != 20 || response.Total != 3 || response.TotalPages != 1 {
		t.Fatalf("unexpected search metadata: %#v", response)
	}
	if len(response.Items) != 3 {
		t.Fatalf("expected 3 matching products, got %#v", response.Items)
	}

	expected := []struct {
		Name string
		SKU  string
	}{
		{Name: "alpha desk", SKU: "WORK-001"},
		{Name: "Beta Desk Shelf", SKU: "SHELF-002"},
		{Name: "Omega Lamp", SKU: "DESK-LIGHT-003"},
	}

	for i, item := range response.Items {
		if item.Name != expected[i].Name || item.SKU != expected[i].SKU {
			t.Fatalf("unexpected search item at index %d: %#v", i, item)
		}
	}
}

func TestSearchProductsReturnsEmptyListWhenThereAreNoMatches(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/search/products?q=missing", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when search has no matches, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response dto.ProductListResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode no-match search response: %v", err)
	}

	if response.Total != 0 || response.TotalPages != 0 || len(response.Items) != 0 {
		t.Fatalf("expected empty search response, got %#v", response)
	}
}

func TestSearchProductsRequiresNonEmptyQuery(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	testCases := []struct {
		name string
		path string
	}{
		{name: "missing q", path: "/v1/search/products"},
		{name: "blank q", path: "/v1/search/products?q=%20%20"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			recorder := performRequest(t, router, http.MethodGet, tt.path, "")
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("expected 400 for invalid search query, got %d with body %s", recorder.Code, recorder.Body.String())
			}

			var response commonapi.ErrorResponse
			if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
				t.Fatalf("failed to decode search error response: %v", err)
			}

			if response.Error.Code != "VALIDATION_ERROR" {
				t.Fatalf("expected validation error code, got %q", response.Error.Code)
			}
			if response.Path != "/v1/search/products" {
				t.Fatalf("expected response path /v1/search/products, got %q", response.Path)
			}
			if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "q" {
				t.Fatalf("expected q validation detail, got %#v", response.Error.Details)
			}
		})
	}
}
