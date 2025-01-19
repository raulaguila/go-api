package service

import (
	"context"

	"github.com/raulaguila/packhub"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

func NewProductService(r domain.ProductRepository) domain.ProductService {
	return &productService{
		productRepository: r,
	}
}

type productService struct {
	productRepository domain.ProductRepository
}

func (s *productService) GenerateProductOutputDTO(p *domain.Product) *dto.ProductOutputDTO {
	return &dto.ProductOutputDTO{
		ID:   &p.ID,
		Name: &p.Name,
	}
}

func (s *productService) GetProductByID(ctx context.Context, id uint) (*dto.ProductOutputDTO, error) {
	p := &domain.Product{Base: domain.Base{ID: id}}
	if err := s.productRepository.GetProduct(ctx, p); err != nil {
		return nil, err
	}

	return s.GenerateProductOutputDTO(p), nil
}

func (s *productService) GetProducts(ctx context.Context, f *filter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error) {
	products, err := s.productRepository.GetProducts(ctx, f)
	if err != nil {
		return nil, err
	}

	count, err := s.productRepository.CountProducts(ctx, f)
	if err != nil {
		return nil, err
	}

	outputProducts := make([]dto.ProductOutputDTO, len(*products))
	for i, product := range *products {
		outputProducts[i] = *s.GenerateProductOutputDTO(&product)
	}

	return &dto.ItemsOutputDTO[dto.ProductOutputDTO]{
		Items: outputProducts,
		Pagination: dto.PaginationDTO{
			CurrentPage: uint(packhub.Max(f.Page, 1)),
			PageSize:    uint(packhub.Max(f.Limit, len(outputProducts))),
			TotalItems:  uint(count),
			TotalPages:  uint(f.CalcPages(count)),
		},
	}, nil
}

func (s *productService) CreateProduct(ctx context.Context, pdto *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	product := new(domain.Product)
	if err := product.Bind(pdto); err != nil {
		return nil, err
	}

	if err := s.productRepository.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return s.GenerateProductOutputDTO(product), nil
}

func (s *productService) UpdateProduct(ctx context.Context, id uint, pdto *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	product := &domain.Product{Base: domain.Base{ID: id}}
	if err := s.productRepository.GetProduct(ctx, product); err != nil {
		return nil, err
	}

	if err := product.Bind(pdto); err != nil {
		return nil, err
	}

	if err := s.productRepository.UpdateProduct(ctx, product); err != nil {
		return nil, err
	}

	return s.GenerateProductOutputDTO(product), nil
}

func (s *productService) DeleteProducts(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return s.productRepository.DeleteProducts(ctx, ids)
}
