package handler

import (
	"fmt"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"github.com/raulaguila/go-api/pkg/filter"
)

func TestProductHandlerSuite(t *testing.T) {
	suite.Run(t, new(ProductHandlerTestSuite))
}

type ProductHandlerTestSuite struct {
	suite.Suite
	filter *filter.Filter
	route  string

	app *fiber.App
}

// SetupSuite function executes before the test suite begins execution
func (s *ProductHandlerTestSuite) SetupSuite() {
	fmt.Println(">>> From SetupSuite")
	s.filter = filter.New("name", "desc")
	s.route = "/product"

	app := fiber.New()
	app.Use(
		fiberi18n.New(&fiberi18n.Config{
			Next: func(c *fiber.Ctx) bool {
				return false
			},
			RootPath:        "./locales",
			AcceptLanguages: []language.Tag{language.AmericanEnglish, language.BrazilianPortuguese},
			DefaultLanguage: language.AmericanEnglish,
			Loader:          &fiberi18n.EmbedLoader{FS: configs.Locales},
		}),
	)
	listProducts := &dto.ItemsOutputDTO[dto.ProductOutputDTO]{
		Items: []dto.ProductOutputDTO{},
		Pagination: dto.PaginationDTO{
			CurrentPage: 0,
			PageSize:    0,
			TotalItems:  0,
			TotalPages:  0,
		},
	}

	service := &mocks.ProductServiceMock{}
	service.On("GetProducts", mock.AnythingOfType("*fasthttp.RequestCtx"), mock.AnythingOfType("*filter.Filter")).Return(listProducts, nil)

	var id1 uint = 1
	name1 := "Product 01"
	service.On("GetProductByID", mock.AnythingOfType("*fasthttp.RequestCtx"), id1).Return(nil, gorm.ErrRecordNotFound)
	service.On("GetProductByID", mock.AnythingOfType("*fasthttp.RequestCtx"), id1).Return(&dto.ProductOutputDTO{ID: &id1, Name: &name1}, nil)
	service.On("UpdateProduct", mock.AnythingOfType("*fasthttp.RequestCtx"), id1, mock.AnythingOfType("*dto.ProductInputDTO")).Return(&dto.ProductOutputDTO{ID: &id1, Name: &name1}, nil)
	service.On("DeleteProducts", mock.AnythingOfType("*fasthttp.RequestCtx"), []uint{id1}).Return(nil)

	var id2 uint = 2
	name2 := "Product 02"
	service.On("GetProductByID", mock.AnythingOfType("*fasthttp.RequestCtx"), id2).Return(&dto.ProductOutputDTO{ID: &id2, Name: &name2}, nil)
	service.On("UpdateProduct", mock.AnythingOfType("*fasthttp.RequestCtx"), id2, mock.AnythingOfType("*dto.ProductInputDTO")).Return(&dto.ProductOutputDTO{ID: &id2, Name: &name2}, nil)
	service.On("DeleteProducts", mock.AnythingOfType("*fasthttp.RequestCtx"), []uint{id2}).Return(gorm.ErrRecordNotFound)

	var id3 uint = 3
	service.On("UpdateProduct", mock.AnythingOfType("*fasthttp.RequestCtx"), id3, mock.AnythingOfType("*dto.ProductInputDTO")).Return(nil, pgutils.ErrDuplicatedKey)

	NewProductHandler(app.Group(s.route), service)
	s.app = app
}

// TearDownSuite function executes after all tests executed
func (s *ProductHandlerTestSuite) TearDownSuite() {
	fmt.Println(">>> From TearDownSuite")
}

// SetupTest function executes before each test case
func (s *ProductHandlerTestSuite) SetupTest() {
	fmt.Println(">>> From SetupTest")
}

// TestCreateProduct test to create a new product
func (s *ProductHandlerTestSuite) TestCreateProduct() {
	fmt.Println(">>> From TestCreateProduct")
}

// TestGetProducts test to list products
func (s *ProductHandlerTestSuite) TestGetProducts() {
	req := httptest.NewRequest(fiber.MethodGet, s.route, nil)
	resp, _ := s.app.Test(req, 100)
	s.Equal(fiber.StatusOK, resp.StatusCode, "Wrong status code.")
	if resp.StatusCode == fiber.StatusOK {
		body, err := io.ReadAll(resp.Body)
		s.NoError(err)
		s.Equal("{\"items\":[],\"pagination\":{\"current_page\":0,\"page_size\":0,\"total_items\":0,\"total_pages\":0}}", string(body))
	}
}

// TestCreateProduct test to create a new product
func (s *ProductHandlerTestSuite) TestGetProductByID() {
	for _, test := range []struct {
		productID    int
		expectedCode int
	}{
		{1, fiber.StatusNotFound},
		{2, fiber.StatusOK},
		{0, fiber.StatusBadRequest},
		{-4, fiber.StatusBadRequest},
	} {
		req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("%v/%v", s.route, test.productID), nil)
		resp, _ := s.app.Test(req, 100)
		s.Equal(test.expectedCode, resp.StatusCode, "Wrong status code.")
		if resp.StatusCode == fiber.StatusOK {
			body, err := io.ReadAll(resp.Body)
			s.NoError(err)
			s.Equal("{\"id\":2,\"name\":\"Product 02\"}", string(body))
		}
	}
}

// TestUpdateProductByID test to update a product
func (s *ProductHandlerTestSuite) TestUpdateProductByID() {
	for _, test := range []struct {
		productID    int
		expectedCode int
		expectedBody string
	}{
		{1, fiber.StatusOK, "{\"id\":1,\"name\":\"Product 01\"}"},
		{2, fiber.StatusOK, "{\"id\":2,\"name\":\"Product 02\"}"},
		{3, fiber.StatusConflict, "{\"code\":409,\"message\":\"Product already registered.\"}"},
		{0, fiber.StatusBadRequest, "{\"code\":400,\"message\":\"Invalid id, please specify valid id.\"}"},
		{-4, fiber.StatusBadRequest, "{\"code\":400,\"message\":\"Invalid id, please specify valid id.\"}"},
	} {
		productDTO := fmt.Sprintf("{\"name\":\"Product 0%v\"}", test.productID)
		fmt.Println(productDTO)

		req := httptest.NewRequest(fiber.MethodPut, fmt.Sprintf("%v/%v", s.route, test.productID), strings.NewReader(productDTO))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := s.app.Test(req, 1000)
		s.Equal(test.expectedCode, resp.StatusCode, "Wrong status code.")
		if resp.StatusCode == fiber.StatusOK {
			body, err := io.ReadAll(resp.Body)
			s.NoError(err)
			s.Equal(test.expectedBody, string(body))
		}
	}
}

// TestDeleteProductByID test to delete a product
func (s *ProductHandlerTestSuite) TestDeleteProductByID() {
	for _, test := range []struct {
		productID    int
		expectedCode int
	}{
		{1, fiber.StatusNoContent},
		{2, fiber.StatusNotFound},
		{-2, fiber.StatusInternalServerError},
	} {
		idsDTO := fmt.Sprintf("{\"ids\": [%v]}", test.productID)
		req := httptest.NewRequest(fiber.MethodDelete, s.route, strings.NewReader(idsDTO))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := s.app.Test(req, 1000)
		s.Equal(test.expectedCode, resp.StatusCode, "Wrong status code.")
	}
}
