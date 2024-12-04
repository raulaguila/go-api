package domain

//
//import (
//	"context"
//	"github.com/raulaguila/go-api/pkg/validator"
//	"testing"
//
//	"github.com/raulaguila/go-api/internal/pkg/dto"
//	"github.com/raulaguila/go-api/pkg/filter"
//	"github.com/stretchr/testify/assert"
//)
//
//type mockValidator struct {
//	err error
//}
//
//func (m *mockValidator) Validate(obj interface{}) error {
//	return m.err
//}
//
//func TestProductTableName(t *testing.T) {
//	product := &Product{}
//	expected := ProductTableName
//	assert.Equal(t, expected, product.TableName(), "they should be equal")
//}
//
//func TestProductBind(t *testing.T) {
//	tests := []struct {
//		name    string
//		input   *dto.ProductInputDTO
//		wantErr bool
//	}{
//		{"Valid input", &dto.ProductInputDTO{Name: strPtr("Product1")}, false},
//		{"Nil input", nil, false},
//		{"Empty name", &dto.ProductInputDTO{Name: strPtr("")}, true},
//		{"Short name", &dto.ProductInputDTO{Name: strPtr("A")}, true},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			product := &Product{}
//			validator.StructValidator = &mockValidator{nil}
//			err := product.Bind(tt.input)
//			if tt.wantErr {
//				assert.Error(t, err, "expected error")
//			} else {
//				assert.NoError(t, err, "expected no error")
//				if tt.input != nil && tt.input.Name != nil {
//					assert.Equal(t, *tt.input.Name, product.Name)
//				}
//			}
//		})
//	}
//}
//
//func TestProductToMap(t *testing.T) {
//	product := &Product{Name: "ProductName"}
//	expected := map[string]any{"name": "ProductName"}
//	assert.Equal(t, expected, product.ToMap(), "they should be equal")
//}
//
//type mockProductRepository struct {
//	products map[uint]*Product
//	err      error
//}
//
//func (repo *mockProductRepository) CountProducts(ctx context.Context, f *filter.Filter) (int64, error) {
//	if repo.err != nil {
//		return 0, repo.err
//	}
//	return int64(len(repo.products)), nil
//}
//
//func (repo *mockProductRepository) GetProductByID(ctx context.Context, id uint) (*Product, error) {
//	if repo.err != nil {
//		return nil, repo.err
//	}
//	return repo.products[id], nil
//}
//
//func (repo *mockProductRepository) GetProducts(ctx context.Context, f *filter.Filter) (*[]Product, error) {
//	if repo.err != nil {
//		return nil, repo.err
//	}
//	result := make([]Product, 0, len(repo.products))
//	for _, p := range repo.products {
//		result = append(result, *p)
//	}
//	return &result, nil
//}
//
//func (repo *mockProductRepository) CreateProduct(ctx context.Context, p *Product) error {
//	if repo.err != nil {
//		return repo.err
//	}
//	repo.products[1] = p
//	return nil
//}
//
//func (repo *mockProductRepository) UpdateProduct(ctx context.Context, p *Product) error {
//	if repo.err != nil {
//		return repo.err
//	}
//	repo.products[1] = p
//	return nil
//}
//
//func (repo *mockProductRepository) DeleteProducts(ctx context.Context, ids []uint) error {
//	if repo.err != nil {
//		return repo.err
//	}
//	for _, id := range ids {
//		delete(repo.products, id)
//	}
//	return nil
//}
//
//func strPtr(s string) *string {
//	return &s
//}
//
//func TestProductRepository(t *testing.T) {
//	repo := &mockProductRepository{
//		products: make(map[uint]*Product),
//	}
//
//	ctx := context.Background()
//	testProduct := &Product{Name: "Test Product"}
//
//	tests := []struct {
//		name     string
//		testFunc func() error
//		wantErr  bool
//	}{
//		{"CreateProduct", func() error { return repo.CreateProduct(ctx, testProduct) }, false},
//		{"CountProducts", func() error { _, err := repo.CountProducts(ctx, nil); return err }, false},
//		{"GetProductByID", func() error { _, err := repo.GetProductByID(ctx, 1); return err }, false},
//		{"GetProducts", func() error { _, err := repo.GetProducts(ctx, nil); return err }, false},
//		{"UpdateProduct", func() error { return repo.UpdateProduct(ctx, testProduct) }, false},
//		{"DeleteProducts", func() error { return repo.DeleteProducts(ctx, []uint{1}) }, false},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			err := tt.testFunc()
//			if tt.wantErr {
//				assert.Error(t, err, "expected error")
//			} else {
//				assert.NoError(t, err, "expected no error")
//			}
//		})
//	}
//}
