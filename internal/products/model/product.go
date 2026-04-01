package model

import commonmodel "go-gin-ecommerce/internal/common/model"

type Product struct {
	ID               uint    `gorm:"primaryKey"`
	Name             string  `gorm:"type:text;not null"`
	StockKeepingUnit string  `gorm:"column:stock_keeping_unit;type:text;not null"`
	Price            float64 `gorm:"type:numeric(12,2);not null"`
	Status           string  `gorm:"type:text;not null"`
	CategoryID       *uint   `gorm:"column:category_id"`
	commonmodel.AuditFields
}

func (Product) TableName() string {
	return "products"
}
