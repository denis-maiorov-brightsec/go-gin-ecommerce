package pagination_test

import (
	"errors"
	"net/url"
	"testing"

	commonapi "go-gin-ecommerce/internal/common/api"
	commonpagination "go-gin-ecommerce/internal/common/pagination"
)

func TestParseUsesDefaults(t *testing.T) {
	t.Parallel()

	params, err := commonpagination.Parse(url.Values{})
	if err != nil {
		t.Fatalf("expected default pagination to parse, got %v", err)
	}

	if params.Page != commonpagination.DefaultPage {
		t.Fatalf("expected default page %d, got %d", commonpagination.DefaultPage, params.Page)
	}

	if params.Limit != commonpagination.DefaultLimit {
		t.Fatalf("expected default limit %d, got %d", commonpagination.DefaultLimit, params.Limit)
	}
}

func TestParseClampsLimitToMax(t *testing.T) {
	t.Parallel()

	params, err := commonpagination.Parse(url.Values{
		"limit": []string{"999"},
	})
	if err != nil {
		t.Fatalf("expected pagination to parse, got %v", err)
	}

	if params.Limit != commonpagination.MaxLimit {
		t.Fatalf("expected limit to be clamped to %d, got %d", commonpagination.MaxLimit, params.Limit)
	}
}

func TestParseRejectsInvalidValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		query url.Values
		field string
	}{
		{
			name:  "invalid page",
			query: url.Values{"page": []string{"0"}},
			field: "page",
		},
		{
			name:  "invalid limit",
			query: url.Values{"limit": []string{"abc"}},
			field: "limit",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := commonpagination.Parse(tt.query)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}

			var apiErr *commonapi.Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected common api error, got %#v", err)
			}

			if apiErr.Code != "VALIDATION_ERROR" {
				t.Fatalf("expected validation error code, got %q", apiErr.Code)
			}

			if len(apiErr.Details) != 1 || apiErr.Details[0].Field != tt.field {
				t.Fatalf("expected validation detail for %s, got %#v", tt.field, apiErr.Details)
			}
		})
	}
}

func TestNewResponseIncludesMetadata(t *testing.T) {
	t.Parallel()

	response := commonpagination.NewResponse([]string{"a", "b"}, commonpagination.Params{
		Page:  2,
		Limit: 2,
	}, 5)

	if response.Page != 2 || response.Limit != 2 || response.Total != 5 || response.TotalPages != 3 {
		t.Fatalf("unexpected pagination response: %#v", response)
	}

	if len(response.Items) != 2 {
		t.Fatalf("expected response items to be preserved, got %#v", response.Items)
	}
}
