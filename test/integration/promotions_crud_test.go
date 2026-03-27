package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	platformauth "go-gin-ecommerce/internal/platform/auth"
	"go-gin-ecommerce/internal/promotions/dto"
	"go-gin-ecommerce/test/integration/testutil"
)

func TestPromotionsCRUD(t *testing.T) {
	router := testutil.NewRouterWithDB(t)
	headers := promotionsAuthHeaders()

	createRecorder := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
		"name": "Spring Sale",
		"code": "SPRING-10",
		"discountType": "percentage",
		"discountValue": 10,
		"startsAt": "2026-04-01T00:00:00Z",
		"endsAt": "2026-04-30T23:59:59Z",
		"status": "active"
	}`, headers)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating promotion, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	created := decodePromotionResponse(t, createRecorder)
	if created.ID == 0 {
		t.Fatalf("expected created promotion to have an id, got %#v", created)
	}
	if created.Name != "Spring Sale" || created.Code != "SPRING-10" || created.DiscountType != "percentage" || created.DiscountValue != 10 || created.Status != "active" {
		t.Fatalf("unexpected create response: %#v", created)
	}
	if created.StartsAt == nil || created.EndsAt == nil {
		t.Fatalf("expected date window to be present, got %#v", created)
	}

	listRecorder := performRequestWithHeaders(t, router, http.MethodGet, "/v1/promotions", "", headers)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when listing promotions, got %d with body %s", listRecorder.Code, listRecorder.Body.String())
	}

	var listed []dto.PromotionResponse
	if err := json.Unmarshal(listRecorder.Body.Bytes(), &listed); err != nil {
		t.Fatalf("failed to decode promotion list: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != created.ID {
		t.Fatalf("unexpected listed promotions: %#v", listed)
	}

	getRecorder := performRequestWithHeaders(t, router, http.MethodGet, "/v1/promotions/1", "", headers)
	if getRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when fetching promotion, got %d with body %s", getRecorder.Code, getRecorder.Body.String())
	}

	updateRecorder := performJSONRequestWithHeaders(t, router, http.MethodPatch, "/v1/promotions/1", `{
		"discountValue": 15,
		"status": "scheduled",
		"endsAt": null
	}`, headers)
	if updateRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when updating promotion, got %d with body %s", updateRecorder.Code, updateRecorder.Body.String())
	}

	updated := decodePromotionResponse(t, updateRecorder)
	if updated.DiscountValue != 15 || updated.Status != "scheduled" {
		t.Fatalf("unexpected updated promotion: %#v", updated)
	}
	if updated.EndsAt != nil {
		t.Fatalf("expected endsAt to be cleared, got %#v", updated.EndsAt)
	}

	deleteRecorder := performRequestWithHeaders(t, router, http.MethodDelete, "/v1/promotions/1", "", headers)
	if deleteRecorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 when deleting promotion, got %d", deleteRecorder.Code)
	}

	missingRecorder := performRequestWithHeaders(t, router, http.MethodGet, "/v1/promotions/1", "", headers)
	if missingRecorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 after deleting promotion, got %d", missingRecorder.Code)
	}
}

func TestCreatePromotionDuplicateCodeReturnsConflictEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)
	headers := promotionsAuthHeaders()

	firstRecorder := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
		"name": "Spring Sale",
		"code": "SPRING-10",
		"discountType": "percentage",
		"discountValue": 10,
		"status": "active"
	}`, headers)
	if firstRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating initial promotion, got %d with body %s", firstRecorder.Code, firstRecorder.Body.String())
	}

	recorder := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
		"name": "Spring Sale Copy",
		"code": "SPRING-10",
		"discountType": "fixed",
		"discountValue": 5,
		"status": "draft"
	}`, headers)
	if recorder.Code != http.StatusConflict {
		t.Fatalf("expected 409 for duplicate code, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "CONFLICT" {
		t.Fatalf("expected CONFLICT error code, got %q", response.Error.Code)
	}
	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "code" {
		t.Fatalf("expected code conflict detail, got %#v", response.Error.Details)
	}
}

func TestCreatePromotionRejectsInvalidDateWindow(t *testing.T) {
	router := testutil.NewRouterWithDB(t)
	headers := promotionsAuthHeaders()

	recorder := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
		"name": "Spring Sale",
		"code": "SPRING-10",
		"discountType": "percentage",
		"discountValue": 10,
		"startsAt": "2026-04-30T00:00:00Z",
		"endsAt": "2026-04-01T00:00:00Z",
		"status": "active"
	}`, headers)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid date window, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "startsAt" {
		t.Fatalf("expected startsAt validation detail, got %#v", response.Error.Details)
	}
}

func TestCreatePromotionRejectsInvalidTimestampFormat(t *testing.T) {
	router := testutil.NewRouterWithDB(t)
	headers := promotionsAuthHeaders()

	recorder := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
		"name": "Spring Sale",
		"code": "SPRING-10",
		"discountType": "percentage",
		"discountValue": 10,
		"startsAt": "not-a-time",
		"status": "active"
	}`, headers)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid timestamp, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "body" {
		t.Fatalf("expected invalid timestamp validation detail, got %#v", response.Error.Details)
	}
}

func TestPatchPromotionRejectsInvalidMergedDateWindow(t *testing.T) {
	router := testutil.NewRouterWithDB(t)
	headers := promotionsAuthHeaders()

	createRecorder := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
		"name": "Spring Sale",
		"code": "SPRING-10",
		"discountType": "percentage",
		"discountValue": 10,
		"endsAt": "2026-04-10T00:00:00Z",
		"status": "active"
	}`, headers)
	if createRecorder.Code != http.StatusCreated {
		t.Fatalf("expected 201 when creating promotion, got %d with body %s", createRecorder.Code, createRecorder.Body.String())
	}

	recorder := performJSONRequestWithHeaders(t, router, http.MethodPatch, "/v1/promotions/1", `{
		"startsAt": "2026-04-11T00:00:00Z"
	}`, headers)
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid merged date window, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
}

func TestDeleteMissingPromotionReturnsNotFoundEnvelope(t *testing.T) {
	router := testutil.NewRouterWithDB(t)
	headers := promotionsAuthHeaders()

	recorder := performRequestWithHeaders(t, router, http.MethodDelete, "/v1/promotions/999", "", headers)
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleting missing promotion, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND error code, got %q", response.Error.Code)
	}
}

func TestPromotionsEndpointsRequireAuthentication(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/promotions", "")
	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for unauthenticated promotions request, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "UNAUTHORIZED" {
		t.Fatalf("expected UNAUTHORIZED error code, got %q", response.Error.Code)
	}
}

func TestPromotionsEndpointsRejectAuthenticatedUsersWithoutPermission(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequestWithHeaders(t, router, http.MethodGet, "/v1/promotions", "", map[string]string{
		"Authorization": "Bearer " + platformauth.StubReadonlyUserToken,
	})
	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for forbidden promotions request, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "FORBIDDEN" {
		t.Fatalf("expected FORBIDDEN error code, got %q", response.Error.Code)
	}
}

func TestNonProtectedEndpointsRemainAccessibleWithoutAuthentication(t *testing.T) {
	router := testutil.NewRouterWithDB(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/health", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 for non-protected health route, got %d with body %s", recorder.Code, recorder.Body.String())
	}
}

func decodePromotionResponse(t *testing.T, recorder *httptest.ResponseRecorder) dto.PromotionResponse {
	t.Helper()

	var response dto.PromotionResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode promotion response: %v", err)
	}

	return response
}

func decodePromotionErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) commonapi.ErrorResponse {
	t.Helper()

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode error response: %v", err)
	}

	return response
}

func TestPromotionResponseDatesRoundTripRFC3339(t *testing.T) {
	timestamp := time.Date(2026, time.April, 1, 12, 30, 0, 0, time.UTC)
	response := dto.PromotionResponse{
		StartsAt: &timestamp,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("failed to marshal response: %v", err)
	}

	if string(data) == "" {
		t.Fatal("expected marshalled response data")
	}
}

func promotionsAuthHeaders() map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + platformauth.StubPromotionsAdminToken,
	}
}
