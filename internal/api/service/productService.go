package service

import (
	"context"
	"github.com/raulaguila/go-api/pkg/utils"

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

func (s *productService) GenerateProductOutputDTO(product *domain.Product) *dto.ProductOutputDTO {
	return &dto.ProductOutputDTO{
		ID:   &product.ID,
		Name: &product.Name,
	}
}

func (s *productService) GetProductByID(ctx context.Context, productID uint) (*dto.ProductOutputDTO, error) {
	product := &domain.Product{Base: domain.Base{ID: productID}}
	if err := s.productRepository.GetProduct(ctx, product); err != nil {
		return nil, err
	}

	return s.GenerateProductOutputDTO(product), nil
}

func (s *productService) GetProducts(ctx context.Context, productFilter *filter.Filter) (*dto.ItemsOutputDTO[dto.ProductOutputDTO], error) {
	products, err := s.productRepository.GetProducts(ctx, productFilter)
	if err != nil {
		return nil, err
	}

	count, err := s.productRepository.CountProducts(ctx, productFilter)
	if err != nil {
		return nil, err
	}

	outputProducts := make([]dto.ProductOutputDTO, 0)
	for _, product := range *products {
		outputProducts = append(outputProducts, *s.GenerateProductOutputDTO(&product))
	}

	return &dto.ItemsOutputDTO[dto.ProductOutputDTO]{
		Items: outputProducts,
		Pagination: dto.PaginationDTO{
			CurrentPage: uint(utils.Max(productFilter.Page, 1)),
			PageSize:    uint(utils.Max(productFilter.Limit, len(outputProducts))),
			TotalItems:  uint(count),
			TotalPages:  uint(productFilter.CalcPages(count)),
		},
	}, nil
}

func (s *productService) CreateProduct(ctx context.Context, data *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	product := new(domain.Product)
	if err := product.Bind(data); err != nil {
		return nil, err
	}

	if err := s.productRepository.CreateProduct(ctx, product); err != nil {
		return nil, err
	}

	return s.GenerateProductOutputDTO(product), nil
}

func (s *productService) UpdateProduct(ctx context.Context, productID uint, data *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	product := &domain.Product{Base: domain.Base{ID: productID}}
	if err := s.productRepository.GetProduct(ctx, product); err != nil {
		return nil, err
	}

	if err := product.Bind(data); err != nil {
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
