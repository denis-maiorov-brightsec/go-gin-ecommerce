package model

import "time"

type Product struct {
	ID         uint      `gorm:"primaryKey"`
	Name       string    `gorm:"type:text;not null"`
	SKU        string    `gorm:"type:text;not null"`
	Price      float64   `gorm:"type:numeric(12,2);not null"`
	Status     string    `gorm:"type:text;not null"`
	CategoryID *uint     `gorm:"column:category_id"`
	CreatedAt  time.Time `gorm:"not null"`
	UpdatedAt  time.Time `gorm:"not null"`
}

func (Product) TableName() string {
	return "products"
}
