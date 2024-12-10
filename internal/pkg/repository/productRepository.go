package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/filter"
)

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

type productRepository struct {
	db *gorm.DB
}

func (s *productRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if filter != nil {
		db = filter.ApplySearchLike(db, "name")
		db = filter.ApplyOrder(db, nil)
	}

	return db
}

func (s *productRepository) CountProducts(ctx context.Context, filter *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	return count, db.Model(new(domain.Product)).Count(&count).Error
}

func (s *productRepository) GetProducts(ctx context.Context, filter *filter.Filter) (*[]domain.Product, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	products := new([]domain.Product)
	return products, db.Find(products).Error
}

func (s *productRepository) GetProductByID(ctx context.Context, productID uint) (*domain.Product, error) {
	product := new(domain.Product)
	return product, s.db.WithContext(ctx).First(product, productID).Error
}

func (s *productRepository) CreateProduct(ctx context.Context, product *domain.Product) error {
	return s.db.WithContext(ctx).Create(product).Error
}

func (s *productRepository) UpdateProduct(ctx context.Context, product *domain.Product) error {
	return s.db.WithContext(ctx).Model(product).Updates(product.ToMap()).Error
}

func (s *productRepository) DeleteProducts(ctx context.Context, toDelete []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(new(domain.Product), toDelete).Error
	})
}
