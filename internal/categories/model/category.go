package model

import "time"

type Category struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"type:text;not null"`
	Slug        string    `gorm:"type:text;not null;uniqueIndex"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

func (Category) TableName() string {
	return "categories"
}
