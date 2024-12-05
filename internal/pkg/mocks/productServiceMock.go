package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

type ProductServiceMock struct {
	mock.Mock
}

func (s *ProductServiceMock) GenerateProductOutputDTO(product *domain.Product) *dto.ProductOutputDTO {
	return &dto.ProductOutputDTO{
		ID:   &product.ID,
		Name: &product.Name,
	}
}

func (s *ProductServiceMock) GetProductByID(ctx context.Context, productID uint) (*dto.ProductOutputDTO, error) {
	ret := s.Called(ctx, productID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*dto.ProductOutputDTO), ret.Error(1)
}

func (s *ProductServiceMock) GetProducts(ctx context.Context, f *filter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error) {
	ret := s.Called(ctx, f)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*dto.ItemsOutputDTO[dto.ProductOutputDTO]), ret.Error(1)
}

func (s *ProductServiceMock) CreateProduct(ctx context.Context, productDTO *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	ret := s.Called(ctx, productDTO)
	return ret.Get(0).(*dto.ProductOutputDTO), ret.Error(1)
}

func (s *ProductServiceMock) UpdateProduct(ctx context.Context, productID uint, productDTO *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	ret := s.Called(ctx, productID, productDTO)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*dto.ProductOutputDTO), ret.Error(1)
}

func (s *ProductServiceMock) DeleteProducts(ctx context.Context, productIDs []uint) error {
	ret := s.Called(ctx, productIDs)
	return ret.Error(0)
}
