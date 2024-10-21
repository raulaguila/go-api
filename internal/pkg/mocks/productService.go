package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
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
	return nil, nil
}

// GetProducts(context.Context, *filter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error)
// CreateProduct(context.Context, *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
// UpdateProduct(context.Context, uint, *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
// DeleteProducts(context.Context, []uint) error
