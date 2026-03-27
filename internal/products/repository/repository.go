package repository

import (
	"context"
	"errors"

	"go-gin-ecommerce/internal/products/model"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id uint) (model.Product, error)
	Create(ctx context.Context, product *model.Product) error
	Update(ctx context.Context, product *model.Product) error
	Delete(ctx context.Context, id uint) error
}

var ErrNotFound = errors.New("product not found")

type GormRepository struct {
	db *gorm.DB
}

func New(db *gorm.DB) *GormRepository {
	return &GormRepository{db: db}
}

func (r *GormRepository) List(ctx context.Context) ([]model.Product, error) {
	var products []model.Product
	if err := r.db.WithContext(ctx).Order("id ASC").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func (r *GormRepository) GetByID(ctx context.Context, id uint) (model.Product, error) {
	var product model.Product
	err := r.db.WithContext(ctx).First(&product, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return model.Product{}, ErrNotFound
	}
	if err != nil {
		return model.Product{}, err
	}

	return product, nil
}

func (r *GormRepository) Create(ctx context.Context, product *model.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *GormRepository) Update(ctx context.Context, product *model.Product) error {
	tx := r.db.WithContext(ctx).Model(&model.Product{}).Where("id = ?", product.ID).Updates(map[string]any{
		"name":        product.Name,
		"sku":         product.SKU,
		"price":       product.Price,
		"status":      product.Status,
		"category_id": product.CategoryID,
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
	tx := r.db.WithContext(ctx).Delete(&model.Product{}, id)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return ErrNotFound
	}

	return nil
}
