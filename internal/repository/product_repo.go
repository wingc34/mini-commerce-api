package repository

import (
	"github.com/wingc34/mini-commerce-api/internal/model"
	"gorm.io/gorm"
)

// ProductRepository defines the database operations for products.
// Defining this as an interface (not a concrete struct) lets us inject
// a mock implementation in tests without touching a real database.
type ProductRepository interface {
	FindByID(id string) (*model.Product, error)
	FindAll(page int, pageSize int) ([]model.Product, int64, error)
	FindRecommended(limit int) ([]model.Product, error)
	CheckStock(productID string, attrsJSON string) (bool, error)
}

// productRepository is the GORM implementation of ProductRepository.
// Lowercase name keeps it private - external packages can only use
// the ProductRepository interface, not this concrete type directly.
type productRepository struct {
	db *gorm.DB
}

// NewProductRepository returns a ProductRepository backed by the given GORM connection.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

// FindByID returns a single product with its SKUs preloaded.
func (r *productRepository) FindByID(id string) (*model.Product, error) {
	var product model.Product
	err := r.db.Preload("SKUs").First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// FindAll returns a paginated list of products along with the total count,
// used by the handler to calculate total pages for the frontend.
func (r *productRepository) FindAll(page int, pageSize int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	offset := (page - 1) * pageSize

	if err := r.db.Model(&model.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("SKUs").
		Offset(offset).
		Limit(pageSize).
		Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// FindRecommended returns a random sample of products for the homepage.
func (r *productRepository) FindRecommended(limit int) ([]model.Product, error) {
	var products []model.Product
	err := r.db.Preload("SKUs").
		Order("RANDOM()").
		Limit(limit).
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) CheckStock(productID string, attrsJSON string) (bool, error) {
	var count int64
	err := r.db.Model(&model.SKU{}).
		Where("product_id = ? AND attributes = ?::jsonb AND stock > 0", productID, attrsJSON).
		Count(&count).Error
	return count > 0, err
}
