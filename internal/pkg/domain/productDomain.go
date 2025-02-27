package domain

import (
	"context"

	"github.com/raulaguila/packhub"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgfilter"
	"github.com/raulaguila/go-api/pkg/validator"
)

const ProductTableName string = "product"

type (
	Product struct {
		Base
		Name string `gorm:"column:name;type:varchar(100);unique;index;not null;" validate:"required,min=2"`
	}

	ProductRepository interface {
		CountProducts(ctx context.Context, f *pgfilter.Filter) (int64, error)
		GetProducts(ctx context.Context, f *pgfilter.Filter) (*[]Product, error)
		GetProduct(ctx context.Context, p *Product) error
		CreateProduct(ctx context.Context, p *Product) error
		UpdateProduct(ctx context.Context, p *Product) error
		DeleteProducts(ctx context.Context, ids []uint) error
	}

	ProductService interface {
		GenerateProductOutputDTO(*Product) *dto.ProductOutputDTO
		GetProducts(ctx context.Context, f *pgfilter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error)
		GetProductByID(ctx context.Context, id uint) (*dto.ProductOutputDTO, error)
		CreateProduct(ctx context.Context, pdto *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
		UpdateProduct(ctx context.Context, id uint, pdto *dto.ProductInputDTO) (*dto.ProductOutputDTO, error)
		DeleteProducts(ctx context.Context, ids []uint) error
	}
)

func (s *Product) TableName() string { return ProductTableName }

func (s *Product) Bind(p *dto.ProductInputDTO) error {
	if p != nil {
		s.Name = packhub.PointerValue(p.Name, s.Name)
	}

	return validator.StructValidator.Validate(s)
}

func (s *Product) ToMap() map[string]any {
	return map[string]any{
		"name": s.Name,
	}
}
