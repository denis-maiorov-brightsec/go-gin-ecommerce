package model

import commonmodel "go-gin-ecommerce/internal/common/model"

type Order struct {
	ID          uint        `gorm:"primaryKey"`
	Status      string      `gorm:"type:text;not null"`
	CustomerID  uint        `gorm:"column:customer_id;not null"`
	TotalAmount float64     `gorm:"column:total_amount;type:numeric(12,2);not null"`
	Items       []OrderItem `gorm:"foreignKey:OrderID"`
	commonmodel.AuditFields
}

func (Order) TableName() string {
	return "orders"
}
