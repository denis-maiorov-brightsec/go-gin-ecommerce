package model

import commonmodel "go-gin-ecommerce/internal/common/model"

type Category struct {
	ID          uint    `gorm:"primaryKey"`
	Name        string  `gorm:"type:text;not null"`
	Slug        string  `gorm:"type:text;not null;uniqueIndex"`
	Description *string `gorm:"type:text"`
	commonmodel.AuditFields
}

func (Category) TableName() string {
	return "categories"
}
