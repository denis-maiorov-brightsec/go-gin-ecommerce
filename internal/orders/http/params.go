package http

import (
	"strconv"

	commonapi "go-gin-ecommerce/internal/common/api"
)

func ParseOrderID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		return 0, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "id",
			Constraints: []string{"id must be a positive integer"},
		}})
	}

	if strconv.IntSize == 32 && id > uint64(^uint32(0)) {
		return 0, commonapi.NewValidationError([]commonapi.ErrorDetail{{
			Field:       "id",
			Constraints: []string{"id must be a positive integer"},
		}})
	}

	return uint(id), nil
}
