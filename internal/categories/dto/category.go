package dto

import (
	"time"

	"go-gin-ecommerce/internal/categories/model"
)

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=1"`
	Slug        string  `json:"slug" binding:"required,min=1"`
	Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1"`
	Slug        *string `json:"slug" binding:"omitempty,min=1"`
	Description *string `json:"description"`
}

type CategoryResponse struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description *string   `json:"description,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func NewCategoryResponse(category model.Category) CategoryResponse {
	return CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

func NewCategoryResponses(categories []model.Category) []CategoryResponse {
	responses := make([]CategoryResponse, 0, len(categories))
	for _, category := range categories {
		responses = append(responses, NewCategoryResponse(category))
	}

	return responses
}
