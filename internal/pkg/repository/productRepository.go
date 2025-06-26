package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/pgfilter"
)

func NewProductRepository(db *gorm.DB) domain.ProductRepository {
	return &productRepository{
		db: db,
	}
}

type productRepository struct {
	db *gorm.DB
}

func (s *productRepository) applyFilter(ctx context.Context, f *pgfilter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if f != nil {
		if f.ID != nil {
			db = db.Where("id = ?", *f.ID)
		}

		if where := f.ApplySearchLike("name"); where != "" {
			db = db.Where(where)
		}
		db = db.Order(f.ApplyOrder(nil))
	}

	return db
}

func (s *productRepository) CountProducts(ctx context.Context, f *pgfilter.Filter) (int64, error) {
	var count int64
	return count, s.applyFilter(ctx, f).Model(new(domain.Product)).Count(&count).Error
}

func (s *productRepository) GetProducts(ctx context.Context, f *pgfilter.Filter) (*[]domain.Product, error) {
	db := s.applyFilter(ctx, f)
	if f != nil {
		if ok, offset, limit := f.ApplyPagination(); ok {
			db = db.Offset(offset).Limit(limit)
		}
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

func (s *productRepository) UpdateProduct(ctx context.Context, p *domain.Product) error {
	return s.db.WithContext(ctx).Model(p).Updates(p.ToMap()).Error
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
