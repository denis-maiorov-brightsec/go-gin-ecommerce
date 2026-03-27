package pagination

import (
	"fmt"
	"math"
	"net/url"
	"strconv"

	commonapi "go-gin-ecommerce/internal/common/api"
)

const (
	DefaultPage  = 1
	DefaultLimit = 20
	MaxLimit     = 100
)

type Params struct {
	Page  int
	Limit int
}

type Response[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"totalPages"`
}

func Parse(values url.Values) (Params, error) {
	page, err := parsePositiveInt(values.Get("page"), "page", DefaultPage)
	if err != nil {
		return Params{}, err
	}

	limit, err := parsePositiveInt(values.Get("limit"), "limit", DefaultLimit)
	if err != nil {
		return Params{}, err
	}

	if limit > MaxLimit {
		limit = MaxLimit
	}

	return Params{
		Page:  page,
		Limit: limit,
	}, nil
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}

func NewResponse[T any](items []T, params Params, total int64) Response[T] {
	return Response[T]{
		Items:      items,
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages(total, params.Limit),
	}
}

func parsePositiveInt(raw string, field string, fallback int) (int, error) {
	if raw == "" {
		return fallback, nil
	}

	value, err := strconv.Atoi(raw)
	if err != nil || value <= 0 {
		return 0, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       field,
			Constraints: []string{fmt.Sprintf("%s must be a positive integer", field)},
		}})
	}

	return value, nil
}

func totalPages(total int64, limit int) int {
	if total == 0 {
		return 0
	}

	return int(math.Ceil(float64(total) / float64(limit)))
}
