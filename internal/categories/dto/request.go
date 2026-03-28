package dto

import (
	"bytes"
	"encoding/json"
)

type CreateCategoryRequest struct {
	Name        string  `json:"name" binding:"required,min=1"`
	Slug        string  `json:"slug" binding:"required,min=1"`
	Description *string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        *string        `json:"name" binding:"omitempty,min=1"`
	Slug        *string        `json:"slug" binding:"omitempty,min=1"`
	Description OptionalString `json:"description"`
}

type OptionalString struct {
	Set   bool
	Null  bool
	Value string
}

func (o *OptionalString) UnmarshalJSON(data []byte) error {
	o.Set = true

	if bytes.Equal(bytes.TrimSpace(data), []byte("null")) {
		o.Null = true
		o.Value = ""
		return nil
	}

	o.Null = false
	return json.Unmarshal(data, &o.Value)
}
