package service

//import (
//	"context"
//	"errors"
//	"reflect"
//	"testing"
//
//	"github.com/raulaguila/go-api/internal/pkg/domain"
//	"github.com/raulaguila/go-api/internal/pkg/dto"
//	"github.com/raulaguila/go-api/internal/pkg/_mocks"
//	"github.com/raulaguila/go-api/pkg/filter"
//	"github.com/raulaguila/go-api/pkg/utils"
//)

//func TestGenerateProductOutputDTO(t *testing.T) {
//	service := &productService{}
//	product := &domain.Product{
//		Base: domain.Base{ID: 1},
//		Name: "Test Product",
//	}
//
//	result := service.GenerateProductOutputDTO(product)
//	expected := &dto.ProductOutputDTO{
//		ID:   &product.ID,
//		Name: &product.Name,
//	}
//
//	if !reflect.DeepEqual(result, expected) {
//		t.Fatalf("expected %v, got %v", expected, result)
//	}
//}
//
//func TestGetProductByID(t *testing.T) {
//	tests := []struct {
//		name     string
//		id       uint
//		repo     *_mocks.ProductRepositoryMock
//		expected *dto.ProductOutputDTO
//		err      error
//	}{
//		{
//			name: "product exists",
//			id:   1,
//			repo: &_mocks.ProductRepositoryMock{
//				products: []*domain.Product{{ID: 1, Name: "Test Product"}},
//			},
//			expected: &dto.ProductOutputDTO{ID: utils.Pointer(uint(1)), Name: utils.Pointer("Test Product")},
//			err:      nil,
//		},
//		{
//			name: "product not found",
//			id:   2,
//			repo: &mockProductRepository{
//				products: []*domain.Product{{ID: 1, Name: "Test Product"}},
//				err:      errors.New("not found"),
//			},
//			expected: nil,
//			err:      errors.New("not found"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			service := &productService{productRepository: tt.repo}
//
//			result, err := service.GetProductByID(context.Background(), tt.id)
//
//			if !reflect.DeepEqual(result, tt.expected) || (err != nil && err.Error() != tt.err.Error()) {
//				t.Fatalf("expected %v, got %v, error %v", tt.expected, result, err)
//			}
//		})
//	}
//}
//
//func TestGetProducts(t *testing.T) {
//	tests := []struct {
//		name     string
//		filter   *filter.Filter
//		repo     *mockProductRepository
//		expected *dto.ItemsOutputDTO[dto.ProductOutputDTO]
//		err      error
//	}{
//		{
//			name:   "products found",
//			filter: &filter.Filter{Page: 1, Limit: 10},
//			repo: &mockProductRepository{
//				products: []*domain.Product{{ID: 1, Name: "Test Product"}},
//				count:    1,
//			},
//			expected: &dto.ItemsOutputDTO[dto.ProductOutputDTO]{
//				Items: []dto.ProductOutputDTO{{ID: utils.UintPointer(1), Name: utils.StringPointer("Test Product")}},
//				Pagination: dto.PaginationDTO{
//					CurrentPage: 1,
//					PageSize:    10,
//					TotalItems:  1,
//					TotalPages:  1,
//				},
//			},
//			err: nil,
//		},
//		{
//			name:   "no products found",
//			filter: &filter.Filter{Page: 1, Limit: 10},
//			repo: &mockProductRepository{
//				products: nil,
//				count:    0,
//			},
//			expected: &dto.ItemsOutputDTO[dto.ProductOutputDTO]{
//				Items: nil,
//				Pagination: dto.PaginationDTO{
//					CurrentPage: 1,
//					PageSize:    10,
//					TotalItems:  0,
//					TotalPages:  1,
//				},
//			},
//			err: nil,
//		},
//		{
//			name:   "error fetching products",
//			filter: &filter.Filter{Page: 1, Limit: 10},
//			repo: &mockProductRepository{
//				err: errors.New("database error"),
//			},
//			expected: nil,
//			err:      errors.New("database error"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			service := &productService{productRepository: tt.repo}
//
//			result, err := service.GetProducts(context.Background(), tt.filter)
//
//			if !reflect.DeepEqual(result, tt.expected) || (err != nil && err.Error() != tt.err.Error()) {
//				t.Fatalf("expected %v, got %v, error %v", tt.expected, result, err)
//			}
//		})
//	}
//}
//
//func TestCreateProduct(t *testing.T) {
//	tests := []struct {
//		name     string
//		input    *dto.ProductInputDTO
//		repo     *mockProductRepository
//		expected *dto.ProductOutputDTO
//		err      error
//	}{
//		{
//			name:  "product created successfully",
//			input: &dto.ProductInputDTO{Name: "New Product"},
//			repo:  &mockProductRepository{},
//			expected: &dto.ProductOutputDTO{
//				ID:   utils.UintPointer(0),
//				Name: utils.StringPointer("New Product"),
//			},
//			err: nil,
//		},
//		{
//			name:  "error creating product",
//			input: &dto.ProductInputDTO{Name: "New Product"},
//			repo: &mockProductRepository{
//				err: errors.New("database error"),
//			},
//			expected: nil,
//			err:      errors.New("database error"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			service := &productService{productRepository: tt.repo}
//
//			result, err := service.CreateProduct(context.Background(), tt.input)
//
//			if !reflect.DeepEqual(result, tt.expected) || (err != nil && err.Error() != tt.err.Error()) {
//				t.Fatalf("expected %v, got %v, error %v", tt.expected, result, err)
//			}
//		})
//	}
//}
//
//func TestUpdateProduct(t *testing.T) {
//	tests := []struct {
//		name     string
//		id       uint
//		input    *dto.ProductInputDTO
//		repo     *mockProductRepository
//		expected *dto.ProductOutputDTO
//		err      error
//	}{
//		{
//			name:  "product updated successfully",
//			id:    1,
//			input: &dto.ProductInputDTO{Name: "Updated Product"},
//			repo: &mockProductRepository{
//				products: []*domain.Product{{ID: 1, Name: "Old Product"}},
//			},
//			expected: &dto.ProductOutputDTO{
//				ID:   utils.UintPointer(1),
//				Name: utils.StringPointer("Updated Product"),
//			},
//			err: nil,
//		},
//		{
//			name:  "product not found",
//			id:    2,
//			input: &dto.ProductInputDTO{Name: "Updated Product"},
//			repo: &mockProductRepository{
//				products: []*domain.Product{{ID: 1, Name: "Old Product"}},
//				err:      errors.New("not found"),
//			},
//			expected: nil,
//			err:      errors.New("not found"),
//		},
//		{
//			name:  "error updating product",
//			id:    1,
//			input: &dto.ProductInputDTO{Name: "Updated Product"},
//			repo: &mockProductRepository{
//				products: []*domain.Product{{ID: 1, Name: "Old Product"}},
//				err:      errors.New("database error"),
//			},
//			expected: nil,
//			err:      errors.New("database error"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			service := &productService{productRepository: tt.repo}
//
//			result, err := service.UpdateProduct(context.Background(), tt.id, tt.input)
//
//			if !reflect.DeepEqual(result, tt.expected) || (err != nil && err.Error() != tt.err.Error()) {
//				t.Fatalf("expected %v, got %v, error %v", tt.expected, result, err)
//			}
//		})
//	}
//}
//
//func TestDeleteProducts(t *testing.T) {
//	tests := []struct {
//		name string
//		ids  []uint
//		repo *mockProductRepository
//		err  error
//	}{
//		{
//			name: "delete products",
//			ids:  []uint{1, 2},
//			repo: &mockProductRepository{},
//			err:  nil,
//		},
//		{
//			name: "empty ids list",
//			ids:  nil,
//			repo: &mockProductRepository{},
//			err:  nil,
//		},
//		{
//			name: "error deleting products",
//			ids:  []uint{1, 2},
//			repo: &mockProductRepository{
//				err: errors.New("database error"),
//			},
//			err: errors.New("database error"),
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			service := &productService{productRepository: tt.repo}
//
//			err := service.DeleteProducts(context.Background(), tt.ids)
//
//			if err != nil && err.Error() != tt.err.Error() {
//				t.Fatalf("expected error %v, got %v", tt.err, err)
//			}
//		})
//	}
//}
