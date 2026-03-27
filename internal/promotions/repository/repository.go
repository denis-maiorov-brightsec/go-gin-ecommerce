package repository

import (
	"context"
	"errors"

	"go-gin-ecommerce/internal/promotions/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context) ([]model.Promotion, error)
	GetByID(ctx context.Context, id uint) (model.Promotion, error)
	Create(ctx context.Context, promotion *model.Promotion) error
	Update(ctx context.Context, promotion *model.Promotion) error
	Delete(ctx context.Context, id uint) error
}

var ErrNotFound = errors.New("promotion not found")

type GormRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context) ([]model.Promotion, error) {
	var promotions []model.Promotion
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&promotions).Error; err != nil {
		return nil, err
	}

	return promotions, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (model.Promotion, error) {
	var promotion model.Promotion
	err := r.db.WithContext(ctx).First(&promotion, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Promotion{}, ErrNotFound
	}
	if err != nil {
		return model.Promotion{}, err
	}

	return promotion, nil
}

func (r *GormRepository) Create(ctx context.Context, promotion *model.Promotion) error {
	return r.db.WithContext(ctx).Create(promotion).Error
}

func (r *GormRepository) Update(ctx context.Context, promotion *model.Promotion) error {
	tx := r.db.WithContext(ctx).Model(&model.Promotion{}).Where("id = ?", promotion.ID).Updates(map[string]any{
		"name":           promotion.Name,
		"code":           promotion.Code,
		"discount_type":  promotion.DiscountType,
		"discount_value": promotion.DiscountValue,
		"starts_at":      promotion.StartsAt,
		"ends_at":        promotion.EndsAt,
		"status":         promotion.Status,
	})
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *GormRepository) Delete(ctx context.Context, id uint) error {
	tx := r.db.WithContext(ctx).Delete(&model.Promotion{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
