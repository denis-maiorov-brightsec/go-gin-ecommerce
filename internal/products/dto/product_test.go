package dto

import "testing"

func TestCreateProductRequestResolvedStockKeepingUnit(t *testing.T) {
	t.Run("uses canonical field", func(t *testing.T) {
		value := "SKU-001"
		request := CreateProductRequest{StockKeepingUnit: &value}

		resolved, err := request.ResolvedStockKeepingUnit()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resolved != value {
			t.Fatalf("expected %q, got %q", value, resolved)
		}
	})

	t.Run("accepts deprecated alias", func(t *testing.T) {
		value := "SKU-001"
		request := CreateProductRequest{SKUAlias: &value}

		resolved, err := request.ResolvedStockKeepingUnit()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resolved != value {
			t.Fatalf("expected %q, got %q", value, resolved)
		}
	})

	t.Run("rejects conflicting values", func(t *testing.T) {
		canonical := "SKU-001"
		alias := "SKU-002"
		request := CreateProductRequest{
			StockKeepingUnit: &canonical,
			SKUAlias:         &alias,
		}

		if _, err := request.ResolvedStockKeepingUnit(); err == nil {
			t.Fatal("expected conflict validation error")
		}
	})
}

func TestUpdateProductRequestResolvedStockKeepingUnit(t *testing.T) {
	t.Run("returns nil when omitted", func(t *testing.T) {
		request := UpdateProductRequest{}

		resolved, err := request.ResolvedStockKeepingUnit()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resolved != nil {
			t.Fatalf("expected nil, got %q", *resolved)
		}
	})

	t.Run("accepts deprecated alias", func(t *testing.T) {
		value := "SKU-003"
		request := UpdateProductRequest{SKUAlias: &value}

		resolved, err := request.ResolvedStockKeepingUnit()
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resolved == nil || *resolved != value {
			t.Fatalf("expected %q, got %#v", value, resolved)
		}
	})

	t.Run("rejects conflicting values", func(t *testing.T) {
		canonical := "SKU-001"
		alias := "SKU-002"
		request := UpdateProductRequest{
			StockKeepingUnit: &canonical,
			SKUAlias:         &alias,
		}

		if _, err := request.ResolvedStockKeepingUnit(); err == nil {
			t.Fatal("expected conflict validation error")
		}
	})
}
