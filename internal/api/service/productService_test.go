package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/stretchr/testify/suite"
)

var (
	ErrItemNotFound = errors.New("item not found")
	ErrInvalidValue = errors.New("invalid value")
)

func TestProductSuit(t *testing.T) {
	suite.Run(t, new(ProductTestSuite))
}

type ProductTestSuite struct {
	suite.Suite
	ctx    context.Context
	filter *filter.Filter

	service  domain.ProductService
	items    []domain.Product
	newItems []domain.Product
	dtos     []dto.ProductInputDTO
}

func (s *ProductTestSuite) SetupTest() {
	s.ctx = context.Background()
	s.filter = filter.New("name", "desc")
	s.items = []domain.Product{
		{Base: domain.Base{ID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Product 1"},
		{Base: domain.Base{ID: 2, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Product 2"},
		{Base: domain.Base{ID: 3, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Product 3"},
		{Base: domain.Base{ID: 4, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Product 4"},
		{Base: domain.Base{ID: 5, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Product 5"},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "Product 6"},
		{Base: domain.Base{ID: 6, CreatedAt: time.Now(), UpdatedAt: time.Now()}, Name: "."},
	}
	s.newItems = []domain.Product{
		{Name: "Product 4"},
		{Name: "Product 5"},
		{Name: "."},
	}
	s.dtos = []dto.ProductInputDTO{
		{Name: &s.newItems[0].Name},
		{Name: &s.newItems[1].Name},
		{Name: &s.newItems[2].Name},
	}

	var nilFilter *filter.Filter = nil
	var nilProduct *domain.Product = nil

	repo := &mocks.ProductRepositoryMock{}
	repo.On("GetProducts", s.ctx, s.filter).Return(&s.items, nil)
	repo.On("GetProducts", s.ctx, nilFilter).Return(&s.items, nil)

	repo.On("CountProducts", s.ctx, s.filter).Return(int64(len(s.items)), nil)
	repo.On("CountProducts", s.ctx, nilFilter).Return(int64(0), errors.New("error to count items"))

	repo.On("GetProductByID", s.ctx, s.items[0].ID).Return(&s.items[0], nil)
	repo.On("GetProductByID", s.ctx, s.items[1].ID).Return(&s.items[1], nil)
	repo.On("GetProductByID", s.ctx, s.items[3].ID).Return(&s.items[3], nil)
	repo.On("GetProductByID", s.ctx, s.items[4].ID).Return(&s.items[4], nil)
	repo.On("GetProductByID", s.ctx, s.items[5].ID).Return(&s.items[5], nil)
	repo.On("GetProductByID", s.ctx, uint(10)).Return(nil, ErrItemNotFound)

	repo.On("CreateProduct", s.ctx, &s.newItems[0]).Return(nil)
	repo.On("CreateProduct", s.ctx, &s.newItems[1]).Return(nil)
	repo.On("CreateProduct", s.ctx, nilProduct).Return(nil, errors.New("error to create item"))

	repo.On("UpdateProduct", s.ctx, &s.items[3]).Return(nil)
	repo.On("UpdateProduct", s.ctx, &s.items[4]).Return(nil)

	repo.On("DeleteProducts", s.ctx, []uint{}).Return(nil)
	repo.On("DeleteProducts", s.ctx, []uint{1}).Return(nil)
	repo.On("DeleteProducts", s.ctx, []uint{10}).Return(ErrItemNotFound)

	s.service = NewProductService(repo)
}

func (s *ProductTestSuite) TestGetProducts() {
	items, err := s.service.GetProducts(s.ctx, s.filter)

	s.NoError(err)
	s.IsType(&dto.ItemsOutputDTO[dto.ProductOutputDTO]{}, items)
	s.Equal(items.Pagination.TotalItems, uint(len(s.items)))

	_, err = s.service.GetProducts(s.ctx, nil)

	s.Error(err)
}

func (s *ProductTestSuite) TestGetProductByID() {
	item, err := s.service.GetProductByID(s.ctx, s.items[0].ID)

	s.NoError(err)
	s.IsType(&dto.ProductOutputDTO{}, item)
	s.Equal(s.items[0].Name, *item.Name)
	s.Equal(s.items[0].ID, *item.ID)

	item, err = s.service.GetProductByID(s.ctx, s.items[1].ID)

	s.NoError(err)
	s.IsType(&dto.ProductOutputDTO{}, item)
	s.Equal(s.items[1].Name, *item.Name)
	s.Equal(s.items[1].ID, *item.ID)

	item, err = s.service.GetProductByID(s.ctx, uint(10))

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)
}

func (s *ProductTestSuite) TestCreateProduct() {
	item, err := s.service.CreateProduct(s.ctx, &s.dtos[0])

	s.NoError(err)
	s.IsType(&dto.ProductOutputDTO{}, item)
	s.Equal(s.newItems[0].Name, *item.Name)
	s.Equal(s.newItems[0].ID, *item.ID)

	item, err = s.service.CreateProduct(s.ctx, &s.dtos[1])

	s.NoError(err)
	s.IsType(&dto.ProductOutputDTO{}, item)
	s.Equal(s.newItems[1].Name, *item.Name)
	s.Equal(s.newItems[1].ID, *item.ID)

	item, err = s.service.CreateProduct(s.ctx, nil)

	s.Error(err)
	s.Nil(item)
}

func (s *ProductTestSuite) TestUpdateProduct() {
	item, err := s.service.UpdateProduct(s.ctx, s.items[3].ID, &s.dtos[0])

	s.NoError(err)
	s.IsType(&dto.ProductOutputDTO{}, item)
	s.Equal(*s.dtos[0].Name, *item.Name)
	s.Equal(s.items[3].ID, *item.ID)

	item, err = s.service.UpdateProduct(s.ctx, s.items[4].ID, &s.dtos[1])

	s.NoError(err)
	s.IsType(&dto.ProductOutputDTO{}, item)
	s.Equal(*s.dtos[1].Name, *item.Name)
	s.Equal(s.items[4].ID, *item.ID)

	item, err = s.service.UpdateProduct(s.ctx, 10, &s.dtos[1])

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
	s.Nil(item)

	item, err = s.service.UpdateProduct(s.ctx, s.items[5].ID, &s.dtos[2])

	s.Error(err)
	s.Nil(item)
}

func (s *ProductTestSuite) TestDeleteProducts() {
	s.NoError(s.service.DeleteProducts(s.ctx, []uint{}))
	s.NoError(s.service.DeleteProducts(s.ctx, []uint{1}))

	err := s.service.DeleteProducts(s.ctx, []uint{10})

	s.Error(err)
	s.True(errors.Is(err, ErrItemNotFound))
}
