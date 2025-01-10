package domain

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/utils"
	"github.com/raulaguila/go-api/pkg/validator"
)

const ProductTableName string = "product"

type (
	Product struct {
		Base
		Name string `gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	ProductRepository interface {
		CountProducts(context.Context, *filter.Filter) (int64, error)
		GetProduct(context.Context, *Product) error
		GetProducts(context.Context, *filter.Filter) (*[]Product, error)
		CreateProduct(context.Context, *Product) error
		UpdateProduct(ctx context.Context, product *Product, update map[string]any) error
		DeleteProducts(context.Context, []uint) error
	}

	ProductService interface {
		GenerateProductOutputDTO(*Product) *dto.ProductOutputDTO
		GetProductByID(context.Context, uint) (*dto.ProductOutputDTO, error)
		GetProducts(context.Context, *filter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error)
		CreateProduct(context.Context, *dto.ProductInputDTO) error
		UpdateProduct(context.Context, uint, *dto.ProductInputDTO) error
		DeleteProducts(context.Context, []uint) error
	}
)

func (s *Product) TableName() string { return ProductTableName }

func (s *Product) Bind(p *dto.ProductInputDTO) error {
	if p != nil {
		s.Name = utils.PointerValue(p.Name, s.Name)
	}

	return validator.StructValidator.Validate(s)
}

func (s *Product) ToMap() map[string]any {
	return map[string]any{
		"name": s.Name,
	}
}

func (s *Product) GenerateUpdateMap(data *dto.ProductInputDTO) map[string]any {
	mapped := map[string]any{}
	if data != nil {
		if data.Name != nil {
			mapped["name"] = *data.Name
		}
	}

	return mapped
}
