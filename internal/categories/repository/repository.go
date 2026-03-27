package repository

import (
	"context"
	"errors"

	"go-gin-ecommerce/internal/categories/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context) ([]model.Category, error)
	GetByID(ctx context.Context, id uint) (model.Category, error)
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id uint) error
}

var ErrNotFound = errors.New("category not found")

type GormRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context) ([]model.Category, error) {
	var categories []model.Category
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (model.Category, error) {
	var category model.Category
	err := r.db.WithContext(ctx).First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Category{}, ErrNotFound
	}
	if err != nil {
		return model.Category{}, err
	}

	return category, nil
}

func (r *GormRepository) Create(ctx context.Context, category *model.Category) error {
	return r.db.WithContext(ctx).Create(category).Error
}

func (r *GormRepository) Update(ctx context.Context, category *model.Category) error {
	tx := r.db.WithContext(ctx).Model(&model.Category{}).Where("id = ?", category.ID).Updates(map[string]any{
		"name":        category.Name,
		"slug":        category.Slug,
		"description": category.Description,
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
	tx := r.db.WithContext(ctx).Delete(&model.Category{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
