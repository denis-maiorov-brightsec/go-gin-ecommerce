package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-gin-ecommerce/internal/categories/dto"
	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/test/integration/testutil"
)

func TestCategoriesCRUD(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	createRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
		"name": "Lighting",
		"slug": "lighting",
		"description": "Products for lights and fixtures"
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating category, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	created := decodeCategoryResponse(t, createRecorder)
	if created.ID == 0 {
		t.Fatalf("expected created category to have an id, got %#v", created)
	}
	if created.Name != "Lighting" || created.Slug != "lighting" {
		t.Fatalf("unexpected create response: %#v", created)
	}
	if created.Description == nil || *created.Description != "Products for lights and fixtures" {
		t.Fatalf("unexpected description in create response: %#v", created.Description)
	}

	listRecorder := performRequest(t, router, http.MethodGet, "/v1/categories", "")
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when listing categories, got %d", listRecorder.Code)
	}

	var listed []dto.CategoryResponse
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &listed); err != nil {
		t.Fatalf("failed to decode category list: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != created.ID {
		t.Fatalf("unexpected listed categories: %#v", listed)
	}

	getRecorder := performRequest(t, router, http.MethodGet, "/v1/categories/1", "")
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when fetching category, got %d", getRecorder.Code)
	}

	updateRecorder := performJSONRequest(t, router, http.MethodPatch, "/v1/categories/1", `{
		"name": "Office Lighting",
		"slug": "office-lighting"
	}`)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when updating category, got %d with body %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	updated := decodeCategoryResponse(t, updateRecorder)
	if updated.Name != "Office Lighting" || updated.Slug != "office-lighting" {
		t.Fatalf("unexpected updated category: %#v", updated)
	}
	if updated.Description == nil || *updated.Description != "Products for lights and fixtures" {
		t.Fatalf("expected description to remain unchanged, got %#v", updated.Description)
	}

	deleteRecorder := performRequest(t, router, http.MethodDelete, "/v1/categories/1", "")
	if deleteRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 when deleting category, got %d", deleteRecorder.Code)
	}

	missingRecorder := performRequest(t, router, http.MethodGet, "/v1/categories/1", "")
	if missingRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deleting category, got %d", missingRecorder.Code)
	}
}

func TestCreateCategoryValidationErrorUsesGlobalEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
		"name": "",
		"slug": ""
	}`)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid create request, got %d", recorder.Code)
	}

	response := decodeCategoryErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
	if response.Path != "/v1/categories" {
		t.Fatalf("expected response path /v1/categories, got %q", response.Path)
	}
}

func TestCreateCategoryDuplicateSlugReturnsConflictEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	firstRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
		"name": "Lighting",
		"slug": "lighting"
	}`)
	if firstRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating initial category, got %d with body %s", firstRecorder.Code, firstRecorder.Body.String())
	}

	recorder := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
		"name": "More Lighting",
		"slug": "lighting"
	}`)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate slug, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeCategoryErrorResponse(t, recorder)
	if response.Error.Code != "CONFLICT" {
		t.Fatalf("expected CONFLICT error code, got %q", response.Error.Code)
	}
	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "slug" {
		t.Fatalf("expected slug conflict details, got %#v", response.Error.Details)
	}
}

func TestPatchCategoryCanClearDescription(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	createRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
		"name": "Lighting",
		"slug": "lighting",
		"description": "Products for lights and fixtures"
	}`)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating category, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	updateRecorder := performJSONRequest(t, router, http.MethodPatch, "/v1/categories/1", `{
		"description": null
	}`)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when clearing description, got %d with body %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	updated := decodeCategoryResponse(t, updateRecorder)
	if updated.Description != nil {
		t.Fatalf("expected description to be cleared, got %#v", updated.Description)
	}
}

func TestPatchCategoryDuplicateSlugReturnsConflictEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	firstRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
		"name": "Lighting",
		"slug": "lighting"
	}`)
	if firstRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating initial category, got %d with body %s", firstRecorder.Code, firstRecorder.Body.String())
	}

	secondRecorder := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
		"name": "Furniture",
		"slug": "furniture"
	}`)
	if secondRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating second category, got %d with body %s", secondRecorder.Code, secondRecorder.Body.String())
	}

	recorder := performJSONRequest(t, router, http.MethodPatch, "/v1/categories/2", `{
		"slug": "lighting"
	}`)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate slug on patch, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeCategoryErrorResponse(t, recorder)
	if response.Error.Code != "CONFLICT" {
		t.Fatalf("expected CONFLICT error code, got %q", response.Error.Code)
	}
	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "slug" {
		t.Fatalf("expected slug conflict details, got %#v", response.Error.Details)
	}
}

func TestGetMissingCategoryReturnsNotFoundEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/categories/999", "")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing category, got %d", recorder.Code)
	}

	response := decodeCategoryErrorResponse(t, recorder)
	if response.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND error code, got %q", response.Error.Code)
	}
}

func TestDeleteMissingCategoryReturnsNotFoundEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequest(t, router, http.MethodDelete, "/v1/categories/999", "")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleting missing category, got %d", recorder.Code)
	}

	response := decodeCategoryErrorResponse(t, recorder)
	if response.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND error code, got %q", response.Error.Code)
	}
}

func decodeCategoryResponse(t *testing.T, recorder *httptest.ResponseRecorder) dto.CategoryResponse {
	t.Helper()

	var response dto.CategoryResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode category response: %v", err)
	}

	return response
}

func decodeCategoryErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) commonapi.ErrorResponse {
	t.Helper()

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	return response
}
