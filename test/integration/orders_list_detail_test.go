package integration_test

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	"go-gin-ecommerce/internal/http/routes"
	"go-gin-ecommerce/internal/orders/dto"
	"go-gin-ecommerce/test/integration/testutil"
	"gorm.io/gorm"
)

func TestListOrdersSupportsPaginationAndFilters(t *testing.T) {
	router, database := newOrdersTestApp(t)

	firstID := seedOrder(t, database, orderSeed{
		Status:      "pending",
		CustomerID:  101,
		TotalAmount: 55.50,
		CreatedAt:   time.Date(2026, time.January, 10, 9, 0, 0, 0, time.UTC),
		Items: []orderItemSeed{
			{Name: "Desk Lamp", Quantity: 1, UnitPrice: 40.00, LineAmount: 40.00},
			{Name: "Bulb", Quantity: 1, UnitPrice: 15.50, LineAmount: 15.50},
		},
	})
	_ = seedOrder(t, database, orderSeed{
		Status:      "fulfilled",
		CustomerID:  102,
		TotalAmount: 20.00,
		CreatedAt:   time.Date(2026, time.January, 11, 11, 30, 0, 0, time.UTC),
		Items: []orderItemSeed{
			{Name: "Mouse Pad", Quantity: 1, UnitPrice: 20.00, LineAmount: 20.00},
		},
	})
	thirdID := seedOrder(t, database, orderSeed{
		Status:      "pending",
		CustomerID:  103,
		TotalAmount: 99.99,
		CreatedAt:   time.Date(2026, time.January, 12, 14, 45, 0, 0, time.UTC),
		Items: []orderItemSeed{
			{Name: "Keyboard", Quantity: 1, UnitPrice: 99.99, LineAmount: 99.99},
		},
	})

	recorder := performRequest(t, router, http.MethodGet, "/v1/orders?status=pending&from=2026-01-10&to=2026-01-12&page=2&limit=1", "")
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when listing filtered orders, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	var response struct {
		Items      []dto.OrderResponse `json:"items"`
		Page       int                 `json:"page"`
		Limit      int                 `json:"limit"`
		Total      int64               `json:"total"`
		TotalPages int                 `json:"totalPages"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode order list response: %v", err)
	}

	if response.Page != 2 || response.Limit != 1 || response.Total != 2 || response.TotalPages != 2 {
		t.Fatalf("unexpected pagination metadata: %#v", response)
	}
	if len(response.Items) != 1 || response.Items[0].ID != thirdID {
		t.Fatalf("expected only the second pending order on page 2, got %#v", response.Items)
	}
	if response.Items[0].Status != "pending" || response.Items[0].CustomerID != 103 || response.Items[0].TotalAmount != 99.99 {
		t.Fatalf("unexpected order payload: %#v", response.Items[0])
	}
	if len(response.Items[0].Items) != 1 || response.Items[0].Items[0].Name != "Keyboard" {
		t.Fatalf("expected order items to be included, got %#v", response.Items[0].Items)
	}

	detailRecorder := performRequest(t, router, http.MethodGet, fmt.Sprintf("/v1/orders/%d", firstID), "")
	if detailRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200 when fetching order detail, got %d with body %s", detailRecorder.Code, detailRecorder.Body.String())
	}

	detail := decodeOrderResponse(t, detailRecorder)
	if detail.ID != firstID || detail.Status != "pending" || detail.CustomerID != 101 {
		t.Fatalf("unexpected order detail: %#v", detail)
	}
	if len(detail.Items) != 2 {
		t.Fatalf("expected order detail items, got %#v", detail.Items)
	}
}

func TestGetMissingOrderReturnsNotFoundEnvelope(t *testing.T) {
	router, _ := newOrdersTestApp(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/orders/999", "")
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for missing order, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeOrderErrorResponse(t, recorder)
	if response.Error.Code != "NOT_FOUND" {
		t.Fatalf("expected NOT_FOUND error code, got %q", response.Error.Code)
	}
}

func TestListOrdersRejectsInvalidDateFilters(t *testing.T) {
	router, _ := newOrdersTestApp(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/orders?from=2026-13-01", "")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid date filter, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeOrderErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
	if response.Path != "/v1/orders" {
		t.Fatalf("expected response path /v1/orders, got %q", response.Path)
	}
	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "from" {
		t.Fatalf("expected from validation detail, got %#v", response.Error.Details)
	}
}

func TestListOrdersRejectsInvalidDateRange(t *testing.T) {
	router, _ := newOrdersTestApp(t)

	recorder := performRequest(t, router, http.MethodGet, "/v1/orders?from=2026-01-12&to=2026-01-10", "")
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid date range, got %d with body %s", recorder.Code, recorder.Body.String())
	}

	response := decodeOrderErrorResponse(t, recorder)
	if response.Error.Code != "VALIDATION_ERROR" {
		t.Fatalf("expected validation error code, got %q", response.Error.Code)
	}
	if len(response.Error.Details) != 1 || response.Error.Details[0].Field != "from" {
		t.Fatalf("expected from validation detail, got %#v", response.Error.Details)
	}
}

type orderSeed struct {
	Status      string
	CustomerID  uint
	TotalAmount float64
	CreatedAt   time.Time
	Items       []orderItemSeed
}

type orderItemSeed struct {
	ProductID  *uint
	Name       string
	Quantity   int
	UnitPrice  float64
	LineAmount float64
}

func newOrdersTestApp(t *testing.T) (http.Handler, *gorm.DB) {
	t.Helper()

	cfg := testutil.NewTestConfig(t)
	database := testutil.NewTestDatabase(t, cfg)

	return routes.NewWithDB(cfg, slog.Default(), database), database
}

func seedOrder(t *testing.T, database *gorm.DB, order orderSeed) uint {
	t.Helper()

	var orderID uint
	row := database.Raw(
		`INSERT INTO orders (status, customer_id, total_amount, created_at, updated_at)
		 VALUES (?, ?, ?, ?, ?)
		 RETURNING id`,
		order.Status,
		order.CustomerID,
		order.TotalAmount,
		order.CreatedAt,
		order.CreatedAt,
	).Row()
	if err := row.Scan(&orderID); err != nil {
		t.Fatalf("failed to insert order: %v", err)
	}

	for _, item := range order.Items {
		if err := database.Exec(
			`INSERT INTO order_items (order_id, product_id, name, quantity, unit_price, line_amount, created_at, updated_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			orderID,
			item.ProductID,
			item.Name,
			item.Quantity,
			item.UnitPrice,
			item.LineAmount,
			order.CreatedAt,
			order.CreatedAt,
		).Error; err != nil {
			t.Fatalf("failed to insert order item: %v", err)
		}
	}

	return orderID
}

func decodeOrderResponse(t *testing.T, recorder *httptest.ResponseRecorder) dto.OrderResponse {
	t.Helper()

	var response dto.OrderResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode order response: %v", err)
	}

	return response
}

func decodeOrderErrorResponse(t *testing.T, recorder *httptest.ResponseRecorder) commonapi.ErrorResponse {
	t.Helper()

	var response commonapi.ErrorResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode order error response: %v", err)
	}

	return response
}
