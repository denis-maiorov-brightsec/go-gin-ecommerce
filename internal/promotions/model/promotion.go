package model

import "time"

type Promotion struct {
	ID            uint       `gorm:"primaryKey"`
	Name          string     `gorm:"type:text;not null"`
	Code          string     `gorm:"type:text;not null;uniqueIndex"`
	DiscountType  string     `gorm:"column:discount_type;type:text;not null"`
	DiscountValue float64    `gorm:"column:discount_value;type:numeric(12,2);not null"`
	StartsAt      *time.Time `gorm:"column:starts_at"`
	EndsAt        *time.Time `gorm:"column:ends_at"`
	Status        string     `gorm:"type:text;not null"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     time.Time  `gorm:"not null"`
}

func (Promotion) TableName() string {
	return "promotions"
}
