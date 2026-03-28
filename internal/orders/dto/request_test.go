package dto

import (
	"net/url"
	"testing"
	"time"
)

func TestParseOrderIDAcceptsPositiveInteger(t *testing.T) {
	id, err := ParseOrderID("42")
	if err != nil {
		t.Fatalf("expected valid order id, got %v", err)
	}

	if id != 42 {
		t.Fatalf("expected order id 42, got %d", id)
	}
}

func TestParseOrderIDRejectsInvalidValues(t *testing.T) {
	for _, rawID := range []string{"", "0", "-1", "abc"} {
		if _, err := ParseOrderID(rawID); err == nil {
			t.Fatalf("expected %q to be rejected", rawID)
		}
	}
}

func TestParseListOrdersParamsParsesFilters(t *testing.T) {
	params, err := ParseListOrdersParams(url.Values{
		"page":   []string{"2"},
		"limit":  []string{"5"},
		"status": []string{"pending"},
		"from":   []string{"2026-01-10"},
		"to":     []string{"2026-01-12"},
	})
	if err != nil {
		t.Fatalf("expected params to parse, got %v", err)
	}

	if params.Pagination.Page != 2 || params.Pagination.Limit != 5 {
		t.Fatalf("unexpected pagination params: %#v", params.Pagination)
	}
	if params.Status != "pending" {
		t.Fatalf("expected status filter to be pending, got %q", params.Status)
	}
	if params.From == nil || !params.From.Equal(time.Date(2026, time.January, 10, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("unexpected from filter: %#v", params.From)
	}
	if params.To == nil || !params.To.Equal(time.Date(2026, time.January, 13, 0, 0, 0, 0, time.UTC)) {
		t.Fatalf("expected to filter to be exclusive next day, got %#v", params.To)
	}
}

func TestParseListOrdersParamsRejectsInvalidDateRange(t *testing.T) {
	_, err := ParseListOrdersParams(url.Values{
		"from": []string{"2026-01-12"},
		"to":   []string{"2026-01-10"},
	})
	if err == nil {
		t.Fatal("expected invalid range to fail")
	}
}

func TestParseListOrdersParamsRejectsInvalidDateFormat(t *testing.T) {
	_, err := ParseListOrdersParams(url.Values{
		"from": []string{"2026-99-10"},
	})
	if err == nil {
		t.Fatal("expected invalid date format to fail")
	}
}
