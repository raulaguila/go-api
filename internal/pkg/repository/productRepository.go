package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/filter"
)

// NewProductRepository creates a new instance of ProductRepository using the provided gorm.DB connection.
func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

// productRepository is a struct that provides access to product-related database operations using Gorm.
type productRepository struct {
	db *gorm.DB
}

// applyFilter applies search and ordering filters to the database query within the provided context.
func (s *productRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if filter != nil {
		db = filter.ApplySearchLike(db, "name")
		db = filter.ApplyOrder(db, nil)
	}

	return db
}

// CountProducts returns the total count of products in the database based on the specified filter criteria.
func (s *productRepository) CountProducts(ctx context.Context, filter *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	return count, db.Model(new(domain.Product)).Count(&count).Error
}

// GetProducts retrieves a list of products from the database based on the provided filter criteria.
// The filter allows for applying search, order, and pagination to the returned product list.
// It returns a pointer to a slice of Product and an error, if any occurs during the database query.
func (s *productRepository) GetProducts(ctx context.Context, filter *filter.Filter) (*[]domain.Product, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	products := new([]domain.Product)
	return products, db.Find(products).Error
}

// GetProductByID retrieves a product from the database using the specified product ID.
func (s *productRepository) GetProductByID(ctx context.Context, productID uint) (*domain.Product, error) {
	product := new(domain.Product)
	return product, s.db.WithContext(ctx).First(product, productID).Error
}

// CreateProduct inserts a new product record into the database within the specified context. Returns an error if the operation fails.
func (s *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.db.WithContext(ctx).Create(product).Error
}

// UpdateProduct updates an existing product in the database with the provided values in the product parameter.
// It takes a context for request scoping and returns an error if the update operation fails.
func (s *productRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.db.WithContext(ctx).Model(product).Updates(product.ToMap()).Error
}

// DeleteProducts removes products identified by the given IDs from the database within a transaction.
// It uses the provided context for database operations and returns an error if the deletion fails.
func (s *productRepository) DeleteProducts(ctx context.Context, toDelete []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(new(domain.Product), toDelete).Error
	})
}
