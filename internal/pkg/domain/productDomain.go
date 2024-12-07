package domain

import (
	"context"
	"github.com/raulaguila/go-api/pkg/utils"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/validator"
)

// ProductTableName is a constant that specifies the database table name for storing product data.
const ProductTableName string = "product"

// Product is a struct representing a product entity with a base and a unique name field.
type (
	Product struct {
		Base
		Name string `gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	ProductRepository interface {
		CountProducts(context.Context, *filter.Filter) (int64, error)
		GetProductByID(context.Context, uint) (*Product, error)
		GetProducts(context.Context, *filter.Filter) (*[]Product, error)
		CreateProduct(context.Context, *Product) error
		UpdateProduct(context.Context, *Product) error
		DeleteProducts(context.Context, []uint) error
	}

	ProductService interface {
		GenerateProductOutputDTO(*Product) *dto.ProductOutputDTO
		GetProductByID(context.Context, uint) (*dto.ProductOutputDTO, error)
		GetProducts(context.Context, *filter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error)
		CreateProduct(context.Context, *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
		UpdateProduct(context.Context, uint, *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
		DeleteProducts(context.Context, []uint) error
	}
)

// TableName returns the table name associated with the Product struct.
func (s *Product) TableName() string { return ProductTableName }

// Bind updates the Product fields based on the provided ProductInputDTO and validates the updated Product structure.
// If the p.Name is non-nil, it assigns the dereferenced p.Name to the Product's Name field.
// The method returns any validation error encountered during the update process.
func (s *Product) Bind(p *dto.ProductInputDTO) error {
	if p != nil {
		s.Name = utils.PointerValue(p.Name, s.Name)
	}

	return validator.StructValidator.Validate(s)
}

// ToMap converts the Product struct into a map, with the key as the field name and the value as the field's value.
func (s *Product) ToMap() map[string]any {
	return map[string]any{
		"name": s.Name,
	}
}
