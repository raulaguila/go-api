package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/pgfilter"
	"github.com/raulaguila/go-api/pkg/utils"
)

func TestProductRepository_CountProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProductRepository(db)

	tests := []struct {
		name          string
		mockSetup     func()
		filter        *pgfilter.Filter
		expectedCount int64
		expectedError error
	}{
		{
			name: "success_count_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 1"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 2"}).Error)
			},
			filter:        nil,
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "success_count_0",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
			},
			filter:        nil,
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			count, err := repository.CountProducts(context.Background(), tt.filter)

			assert.Equal(t, tt.expectedCount, count)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestProductRepository_GetProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProductRepository(db)

	tests := []struct {
		name          string
		mockSetup     func()
		filter        *pgfilter.Filter
		expectedNames []string
		expectedErr   error
	}{
		{
			name: "success_get_products_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 1"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 2"}).Error)
			},
			filter:        pgfilter.New("name", "asc"),
			expectedNames: []string{"Product 1", "Product 2"},
			expectedErr:   nil,
		},
		{
			name: "success_get_products_0",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
			},
			filter:        pgfilter.New("name", "asc"),
			expectedNames: []string{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			data, err := repository.GetProducts(context.Background(), tt.filter)

			for i, name := range tt.expectedNames {
				assert.Equal(t, name, (*data)[i].Name)
			}
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestProductRepository_GetProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProductRepository(db)

	tests := []struct {
		name         string
		mockSetup    func()
		productInput *domain.Product
		expectedName string
		expectedErr  error
	}{
		{
			name: "success_get_product_1",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 1"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 2"}).Error)
			},
			productInput: &domain.Product{Base: domain.Base{ID: 1}},
			expectedName: "Product 1",
			expectedErr:  nil,
		},
		{
			name: "success_get_product_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 1"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 2"}).Error)
			},
			productInput: &domain.Product{Base: domain.Base{ID: 2}},
			expectedName: "Product 2",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err = repository.GetProduct(context.Background(), tt.productInput)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedName, tt.productInput.Name)
		})
	}
}

func TestProductRepository_CreateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	utils.PanicIfErr(err)
	repository := NewProductRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		input       *domain.Product
		expectedErr error
	}{
		{
			name: "success_create_product",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
			},
			input:       &domain.Product{Name: "Product 1"},
			expectedErr: nil,
		},
		{
			name: "error_duplicated_product",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 1"}).Error)
			},
			input:       &domain.Product{Name: "Product 1"},
			expectedErr: errors.New("UNIQUE constraint failed: product.name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err = repository.CreateProduct(context.Background(), tt.input)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestProductRepository_UpdateProduct(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProductRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		input       *domain.Product
		expectedErr error
	}{
		{
			name: "success_update_product",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 1"}).Error)
			},
			input:       &domain.Product{Base: domain.Base{ID: 1}, Name: "Updated Product 1"},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			assert.Equal(t, tt.expectedErr, repository.UpdateProduct(context.Background(), tt.input))
		})
	}
}

func TestProductRepository_DeleteProducts(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProductRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		toDelete    []uint
		expectedErr error
	}{
		{
			name: "success_delete_products",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Product{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Product{}))
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 1"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 2"}).Error)
				utils.PanicIfErr(db.Create(&domain.Product{Name: "Product 3"}).Error)
			},
			toDelete:    []uint{1, 2, 3},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err = repository.DeleteProducts(context.Background(), tt.toDelete)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
