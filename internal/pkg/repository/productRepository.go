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

func (s *productRepository) applyFilter(ctx context.Context, f *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if f != nil {
		db = f.ApplySearchLike(db, "name")
		db = f.ApplyOrder(db, nil)
	}

	return db
}

func (s *productRepository) CountProducts(ctx context.Context, f *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, f)
	return count, db.Model(new(domain.Product)).Count(&count).Error
}

func (s *productRepository) GetProducts(ctx context.Context, f *filter.Filter) (*[]domain.Product, error) {
	db := s.applyFilter(ctx, f)
	if f != nil {
		db = f.ApplyPagination(db)
	}

	products := new([]domain.Product)
	return products, db.Find(products).Error
}

func (s *productRepository) GetProduct(ctx context.Context, p *domain.Product) error {
	return s.db.WithContext(ctx).Where(p).First(p).Error
}

func (s *productRepository) CreateProduct(ctx context.Context, p *domain.Product) error {
	return s.db.WithContext(ctx).Create(p).Error
}

func (s *productRepository) UpdateProduct(ctx context.Context, p *domain.Product, m map[string]any) error {
	tx := s.db.WithContext(ctx).Table(p.TableName()).Where(p).Updates(m)
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *productRepository) DeleteProducts(ctx context.Context, i []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(new(domain.Product), i)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
