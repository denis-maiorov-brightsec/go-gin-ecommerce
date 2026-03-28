package dto

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	commonapi "go-gin-ecommerce/internal/common/api"
	commonpagination "go-gin-ecommerce/internal/common/pagination"
)

const dateLayout = "2006-01-02"

type ListOrdersParams struct {
	Pagination commonpagination.Params
	Status     string
	From       *time.Time
	To         *time.Time
}

func ParseOrderID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		return 0, invalidOrderIDError()
	}

	if strconv.IntSize == 32 && id > uint64(^uint32(0)) {
		return 0, invalidOrderIDError()
	}

	return uint(id), nil
}

func ParseListOrdersParams(values url.Values) (ListOrdersParams, error) {
	paginationParams, err := commonpagination.Parse(values)
	if err != nil {
		return ListOrdersParams{}, err
	}

	from, err := parseDateFilter(values.Get("from"), "from")
	if err != nil {
		return ListOrdersParams{}, err
	}

	to, err := parseDateFilter(values.Get("to"), "to")
	if err != nil {
		return ListOrdersParams{}, err
	}
	if to != nil {
		toValue := to.Add(24 * time.Hour)
		to = &toValue
	}

	if from != nil && to != nil && from.After(to.Add(-24*time.Hour)) {
		return ListOrdersParams{}, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "from",
			Constraints: []string{"from must be before or equal to to"},
		}})
	}

	return ListOrdersParams{
		Pagination: paginationParams,
		Status:     values.Get("status"),
		From:       from,
		To:         to,
	}, nil
}

func parseDateFilter(raw string, field string) (*time.Time, error) {
	if raw == "" {
		return nil, nil
	}

	value, err := time.Parse(dateLayout, raw)
	if err != nil {
		return nil, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       field,
			Constraints: []string{fmt.Sprintf("%s must be a valid date in %s format", field, dateLayout)},
		}})
	}

	utcValue := value.UTC()
	return &utcValue, nil
}

func invalidOrderIDError() error {
	return commonapi.NewValidationError([]commonapi.ErrorDetail{{
		Field:       "id",
		Constraints: []string{"id must be a positive integer"},
	}})
}
