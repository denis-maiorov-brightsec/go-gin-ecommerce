package repository

import (
	"context"
	"errors"
	"time"

	"go-gin-ecommerce/internal/orders/dto"
	"go-gin-ecommerce/internal/orders/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	List(ctx context.Context, params dto.ListOrdersParams) ([]model.Order, int64, error)
	GetByID(ctx context.Context, id uint) (model.Order, error)
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

func (r *GormRepository) List(ctx context.Context, params dto.ListOrdersParams) ([]model.Order, int64, error) {
	query := r.filteredQuery(ctx, params)

	var total int64
	if err := query.Model(&model.Order{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var orders []model.Order
	if err := query.
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		Order("id ASC").
		Limit(params.Pagination.Limit).
		Offset(params.Pagination.Offset()).
		Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (model.Order, error) {
	var order model.Order
	err := r.db.WithContext(ctx).
		Preload("Items", func(db *gorm.DB) *gorm.DB {
			return db.Order("id ASC")
		}).
		First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Order{}, ErrNotFound
	}
	if err != nil {
		return model.Order{}, err
	}

	return order, nil
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

func (r *GormRepository) filteredQuery(ctx context.Context, params dto.ListOrdersParams) *gorm.DB {
	query := r.db.WithContext(ctx)
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.From != nil {
		query = query.Where("created_at >= ?", *params.From)
	}
	if params.To != nil {
		query = query.Where("created_at < ?", *params.To)
	}

	return query
}
