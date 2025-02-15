package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/raulaguila/packhub"

	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgfilter"
)

func TestProductService_GetProductByID(t *testing.T) {
	mockRepository := new(_mocks.ProductRepositoryMock)
	service := NewProductService(mockRepository)

	tests := []struct {
		name      string
		setupMock func()
		productID uint
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("GetProduct", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			productID: 1,
			wantErr:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetProduct", mock.Anything, mock.Anything).
					Return(gorm.ErrRecordNotFound).
					Once()
			},
			productID: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.GetProductByID(context.Background(), tt.productID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_GetProducts(t *testing.T) {
	mockRepository := new(_mocks.ProductRepositoryMock)
	service := NewProductService(mockRepository)

	tests := []struct {
		name      string
		setupMock func()
		productID uint
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("GetProducts", mock.Anything, mock.Anything).
					Return(&[]domain.Product{{
						Base: domain.Base{ID: 1},
						Name: "Product 01",
					}}, nil).
					Once()
				mockRepository.
					On("CountProducts", mock.Anything, mock.Anything).
					Return(int64(1), nil).
					Once()
			},
			productID: 1,
			wantErr:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetProducts", mock.Anything, mock.Anything).
					Return(nil, gorm.ErrRecordNotFound).
					Once()
			},
			productID: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.GetProducts(context.Background(), &pgfilter.Filter{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	mockRepository := new(_mocks.ProductRepositoryMock)
	service := NewProductService(mockRepository)

	tests := []struct {
		name         string
		setupMock    func()
		productInput *dto.ProductInputDTO
		wantErr      bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("CreateProduct", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("GetProduct", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			productInput: &dto.ProductInputDTO{Name: packhub.Pointer("Product 01")},
			wantErr:      false,
		},
		{
			name: "create error",
			setupMock: func() {
				mockRepository.
					On("CreateProduct", mock.Anything, mock.Anything).
					Return(gorm.ErrDuplicatedKey).
					Once()
			},
			productInput: &dto.ProductInputDTO{Name: packhub.Pointer("Product 01")},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.CreateProduct(context.Background(), tt.productInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_UpdateProduct(t *testing.T) {
	mockRepository := new(_mocks.ProductRepositoryMock)
	service := NewProductService(mockRepository)

	tests := []struct {
		name         string
		setup        func()
		productID    uint
		productInput *dto.ProductInputDTO
		wantErr      bool
	}{
		{
			name: "success",
			setup: func() {
				mockRepository.
					On("GetProduct", mock.Anything, mock.Anything).
					Return(nil).
					Twice()
				mockRepository.
					On("UpdateProduct", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			productID:    1,
			productInput: &dto.ProductInputDTO{Name: packhub.Pointer("Product 01")},
			wantErr:      false,
		},
		{
			name: "create error",
			setup: func() {
				mockRepository.
					On("GetProduct", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("UpdateProduct", mock.Anything, mock.Anything).
					Return(gorm.ErrDuplicatedKey).
					Once()
			},
			productID:    1,
			productInput: &dto.ProductInputDTO{Name: packhub.Pointer("Product 01")},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := service.UpdateProduct(context.Background(), tt.productID, tt.productInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
