package service

import (
	"context"
	"github.com/raulaguila/go-api/pkg/utils"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

// NewProductService creates a new instance of ProductService with the provided ProductRepository implementation.
func NewProductService(r domain.ProductRepository) domain.ProductService {
	return &productService{
		productRepository: r,
	}
}

// productService is a struct implementing the ProductService interface, handling product-related business logic.
type productService struct {
	productRepository domain.ProductRepository
}

// GenerateProductOutputDTO converts a domain.Product to a dto.ProductOutputDTO by mapping fields like ID and Name.
func (s *productService) GenerateProductOutputDTO(product *domain.Product) *dto.ProductOutputDTO {
	return &dto.ProductOutputDTO{
		ID:   &product.ID,
		Name: &product.Name,
	}
}

// GetProductByID retrieves a product by its unique ID and returns a ProductOutputDTO or an error if the product is not found.
func (s *productService) GetProductByID(ctx context.Context, productID uint) (*dto.ProductOutputDTO, error) {
	product, err := s.productRepository.GetProductByID(ctx, productID)
	if err != nil {
		return nil, err
	}

	return s.GenerateProductOutputDTO(product), nil
}

// GetProducts retrieves a list of products based on the provided filter criteria and returns a paginated result.
// It queries the product repository for matching products and counts the total number of items.
// The response includes both the product data and pagination details.
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

// CreateProduct creates a new product in the repository using the provided input data and returns the created product.
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

// UpdateProduct updates an existing product using the provided productID and ProductInputDTO data.
// It retrieves the product, binds the input data, updates the product in the repository, and returns the updated product.
func (s *productService) UpdateProduct(ctx context.Context, productID uint, data *dto.ProductInputDTO) (*dto.ProductOutputDTO, error) {
	product, err := s.productRepository.GetProductByID(ctx, productID)
	if err != nil {
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

// DeleteProducts removes a list of products identified by their IDs.
// If the provided list of IDs is empty, the method returns immediately with no action performed.
// Returns an error if the deletion process encounters any issues, otherwise returns nil.
func (s *productService) DeleteProducts(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return s.productRepository.DeleteProducts(ctx, ids)
}
