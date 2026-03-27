package repository

import (
	"context"
	"errors"
	"time"

	"go-gin-ecommerce/internal/orders/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Cancel(ctx context.Context, id uint, updatedAt time.Time) (model.Order, error)
}

var ErrNotFound = errors.New("order not found")
var ErrInvalidTransition = errors.New("order transition not allowed")

type GormRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Cancel(ctx context.Context, id uint, updatedAt time.Time) (model.Order, error) {
	order := model.Order{}
	result := r.db.WithContext(ctx).
		Model(&model.Order{}).
		Clauses(clause.Returning{}).
		Where("id = ? AND status = ?", id, "pending").
		Updates(map[string]any{
			"status":     "cancelled",
			"updated_at": updatedAt,
		})
	if result.Error != nil {
		return model.Order{}, result.Error
	}
	if result.RowsAffected > 0 {
		if err := r.db.WithContext(ctx).
			Preload("Items", func(db *gorm.DB) *gorm.DB {
				return db.Order("id ASC")
			}).
			First(&order, id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return model.Order{}, ErrNotFound
			}

			return model.Order{}, err
		}

		return order, nil
	}

	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Order{}).Where("id = ?", id).Count(&count).Error; err != nil {
		return model.Order{}, err
	}
	if count == 0 {
		return model.Order{}, ErrNotFound
	}

	return model.Order{}, ErrInvalidTransition
}
