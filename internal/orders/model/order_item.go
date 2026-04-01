package model

import commonmodel "go-gin-ecommerce/internal/common/model"

type OrderItem struct {
	ID         uint    `gorm:"primaryKey"`
	OrderID    uint    `gorm:"column:order_id;not null"`
	ProductID  *uint   `gorm:"column:product_id"`
	Name       string  `gorm:"type:text;not null"`
	Quantity   int     `gorm:"not null"`
	UnitPrice  float64 `gorm:"column:unit_price;type:numeric(12,2);not null"`
	LineAmount float64 `gorm:"column:line_amount;type:numeric(12,2);not null"`
	commonmodel.AuditFields
}

func (OrderItem) TableName() string {
	return "order_items"
}
