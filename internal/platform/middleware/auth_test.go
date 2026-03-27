package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	platformauth "go-gin-ecommerce/internal/platform/auth"
	"go-gin-ecommerce/internal/platform/middleware"

	"github.com/gin-gonic/gin"
)

func TestRequirePermissionRejectsUnauthenticatedRequests(t *testing.T) {
	t.Parallel()

	router := newTestRouter()
	router.GET("/v1/protected", middleware.RequirePermission(platformauth.NewStubAuthenticator(), platformauth.PermissionManagePromotions), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 for missing auth header, got %d", recorder.Code)
	}

	response := decodeErrorResponse(t, recorder)
	assertErrorPayload(t, response, "UNAUTHORIZED", "Authentication required")
}

func TestRequirePermissionRejectsForbiddenRequests(t *testing.T) {
	t.Parallel()

	router := newTestRouter()
	router.GET("/v1/protected", middleware.RequirePermission(platformauth.NewStubAuthenticator(), platformauth.PermissionManagePromotions), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+platformauth.StubReadonlyUserToken)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for missing permission, got %d", recorder.Code)
	}

	response := decodeErrorResponse(t, recorder)
	assertErrorPayload(t, response, "FORBIDDEN", "Forbidden")
}

func TestRequirePermissionAllowsAuthorizedRequests(t *testing.T) {
	t.Parallel()

	router := newTestRouter()
	router.GET("/v1/protected", middleware.RequirePermission(platformauth.NewStubAuthenticator(), platformauth.PermissionManagePromotions), func(c *gin.Context) {
		c.Status(http.StatusNoContent)
	})

	req := httptest.NewRequest(http.MethodGet, "/v1/protected", nil)
	req.Header.Set("Authorization", "Bearer "+platformauth.StubPromotionsAdminToken)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusNoContent {
		t.Fatalf("expected 204 for authorized request, got %d", recorder.Code)
	}
}

func assertErrorPayload(t *testing.T, response commonapi.ErrorResponse, code string, message string) {
	t.Helper()

	if response.Error.Code != code {
		t.Fatalf("expected error code %q, got %q", code, response.Error.Code)
	}

	if response.Error.Message != message {
		t.Fatalf("expected error message %q, got %q", message, response.Error.Message)
	}
}
