package _mocks

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/stretchr/testify/mock"
)

func NewProductRepositoryMock() domain.ProductRepository {
	return new(ProductRepositoryMock)
}

type ProductRepositoryMock struct {
	mock.Mock
}

func (s *ProductRepositoryMock) CountProducts(ctx context.Context, f *filter.Filter) (int64, error) {
	ret := s.Called(ctx, f)
	return ret.Get(0).(int64), ret.Error(1)
}

func (s *ProductRepositoryMock) GetProduct(ctx context.Context, product *domain.Product) error {
	ret := s.Called(ctx, product)
	return ret.Error(0)
}

func (s *ProductRepositoryMock) GetProducts(ctx context.Context, f *filter.Filter) (*[]domain.Product, error) {
	ret := s.Called(ctx, f)
	return ret.Get(0).(*[]domain.Product), ret.Error(1)
}

func (s *ProductRepositoryMock) CreateProduct(ctx context.Context, product *domain.Product) error {
	ret := s.Called(ctx, product)
	return ret.Error(0)
}

func (s *ProductRepositoryMock) UpdateProduct(ctx context.Context, product *domain.Product, update map[string]any) error {
	ret := s.Called(ctx, product, update)
	return ret.Error(0)
}

func (s *ProductRepositoryMock) DeleteProducts(ctx context.Context, productIDs []uint) error {
	ret := s.Called(ctx, productIDs)
	return ret.Error(0)
}
