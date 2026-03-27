package integration_test

import (
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-gin-ecommerce/internal/http/routes"
	"go-gin-ecommerce/test/integration/testutil"
	"gorm.io/gorm"
)

func TestWriteRateLimitAppliesToCRUDWriteRoutes(t *testing.T) {
	t.Run("products create is limited", func(t *testing.T) {
		router := newRateLimitedRouterWithDB(t)

		first := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
			"name": "Rate Limited Product One",
			"sku": "rl-product-1",
			"price": 19.99,
			"status": "active",
			"stock": 4
		}`)
		if first.Code != http.StatusCreated {
			t.Fatalf("expected first product create to succeed, got %d with body %s", first.Code, first.Body.String())
		}

		second := performJSONRequest(t, router, http.MethodPost, "/v1/products", `{
			"name": "Rate Limited Product Two",
			"sku": "rl-product-2",
			"price": 24.99,
			"status": "active",
			"stock": 2
		}`)
		assertRateLimited(t, second)
	})

	t.Run("categories create is limited", func(t *testing.T) {
		router := newRateLimitedRouterWithDB(t)

		first := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
			"name": "Rate Limited Category One",
			"slug": "rate-limited-category-one",
			"description": "first"
		}`)
		if first.Code != http.StatusCreated {
			t.Fatalf("expected first category create to succeed, got %d with body %s", first.Code, first.Body.String())
		}

		second := performJSONRequest(t, router, http.MethodPost, "/v1/categories", `{
			"name": "Rate Limited Category Two",
			"slug": "rate-limited-category-two",
			"description": "second"
		}`)
		assertRateLimited(t, second)
	})

	t.Run("promotions create is limited", func(t *testing.T) {
		router := newRateLimitedRouterWithDB(t)
		headers := promotionsAuthHeaders()

		first := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
			"name": "Rate Limited Promotion One",
			"code": "RATE-LIMIT-ONE",
			"discountType": "percentage",
			"discountValue": 10,
			"status": "active"
		}`, headers)
		if first.Code != http.StatusCreated {
			t.Fatalf("expected first promotion create to succeed, got %d with body %s", first.Code, first.Body.String())
		}

		second := performJSONRequestWithHeaders(t, router, http.MethodPost, "/v1/promotions", `{
			"name": "Rate Limited Promotion Two",
			"code": "RATE-LIMIT-TWO",
			"discountType": "fixed",
			"discountValue": 5,
			"status": "draft"
		}`, headers)
		assertRateLimited(t, second)
	})
}

func TestWriteRateLimitAppliesToOrderCancelStateTransition(t *testing.T) {
	t.Run("cancel route is limited", func(t *testing.T) {
		router, database := newRateLimitedOrdersTestApp(t)

		firstOrderID := seedOrder(t, database, orderSeed{
			Status:      "pending",
			CustomerID:  301,
			TotalAmount: 12.34,
			CreatedAt:   time.Date(2026, time.March, 1, 10, 0, 0, 0, time.UTC),
			Items: []orderItemSeed{
				{Name: "Notebook", Quantity: 1, UnitPrice: 12.34, LineAmount: 12.34},
			},
		})

		first := performRequest(t, router, http.MethodPost, fmt.Sprintf("/v1/orders/%d/cancel", firstOrderID), "")
		if first.Code != http.StatusOK {
			t.Fatalf("expected first order cancel to succeed, got %d with body %s", first.Code, first.Body.String())
		}

		secondOrderID := seedOrder(t, database, orderSeed{
			Status:      "pending",
			CustomerID:  302,
			TotalAmount: 45.67,
			CreatedAt:   time.Date(2026, time.March, 1, 11, 0, 0, 0, time.UTC),
			Items: []orderItemSeed{
				{Name: "Headphones", Quantity: 1, UnitPrice: 45.67, LineAmount: 45.67},
			},
		})

		second := performRequest(t, router, http.MethodPost, fmt.Sprintf("/v1/orders/%d/cancel", secondOrderID), "")
		assertRateLimited(t, second)
	})
}

func TestReadOnlyRoutesRemainUnaffectedByWriteRateLimit(t *testing.T) {
	router := newRateLimitedRouterWithDB(t)

	for attempt := 1; attempt <= 3; attempt++ {
		recorder := performRequest(t, router, http.MethodGet, "/v1/health", "")
		if recorder.Code != http.StatusOK {
			t.Fatalf("expected read-only request %d to remain unaffected, got %d with body %s", attempt, recorder.Code, recorder.Body.String())
		}
	}
}

func newRateLimitedRouterWithDB(t *testing.T) http.Handler {
	t.Helper()

	cfg := testutil.NewTestConfig(t)
	cfg.WriteRateLimitRequests = 1
	cfg.WriteRateLimitWindow = time.Hour

	return testutil.NewRouterWithConfigAndDB(t, cfg)
}

func newRateLimitedOrdersTestApp(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()

	cfg := testutil.NewTestConfig(t)
	cfg.WriteRateLimitRequests = 1
	cfg.WriteRateLimitWindow = time.Hour
	database := testutil.NewTestDatabase(t, cfg)

	return routes.NewWithDB(cfg, slog.Default(), database), database
}

func assertRateLimited(t *testing.T, recorder *httptest.ResponseRecorder) {
	t.Helper()

	if recorder.Code != http.StatusTooManyRequests {
		t.Fatalf("expected 429 once write rate limit is exceeded, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodePromotionErrorResponse(t, recorder)
	if response.Error.Code != "TOO_MANY_REQUESTS" {
		t.Fatalf("expected TOO_MANY_REQUESTS error code, got %q", response.Error.Code)
	}
	if response.Error.Message != "Rate limit exceeded" {
		t.Fatalf("expected rate-limit message, got %q", response.Error.Message)
	}
}
