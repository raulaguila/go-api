package handler

import (
	"errors"
	"os"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/middleware"
	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgutils"
)

func setupProductApp(mockService *_mocks.ProductServiceMock) *fiber.App {
	middleware.MidAccess = middleware.Auth(os.Getenv("ACCESS_TOKEN_PUBLIC"), &_mocks.UserRepositoryMock{})

	app := fiber.New()
	app.Use(fiberi18n.New(&fiberi18n.Config{
		Next: func(c *fiber.Ctx) bool {
			return false
		},
		RootPath:        "./locales",
		AcceptLanguages: []language.Tag{language.AmericanEnglish, language.BrazilianPortuguese},
		DefaultLanguage: language.AmericanEnglish,
		Loader:          &fiberi18n.EmbedLoader{FS: configs.Locales},
	}))
	NewProductHandler(app.Group(""), mockService)

	return app
}

func TestProductHandler_getProducts(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProducts", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.ProductOutputDTO]{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			method:   fiber.MethodGet,
			endpoint: "/",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProducts", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}
	runGeneralHandlerTests(t, tests, setupProductApp(mockService))
}

func TestProductHandler_getProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/1",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProductByID", mock.Anything, uint(1)).Return(&dto.ProductOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			method:   fiber.MethodGet,
			endpoint: "/200",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProductByID", mock.Anything, uint(200)).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodGet,
			endpoint: "/500",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProductByID", mock.Anything, uint(500)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupProductApp(mockService))
}

func TestProductHandler_createProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setupMocks: func() {
				mockService.On("CreateProduct", mock.Anything, mock.Anything).Return(&dto.ProductOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "duplicate product",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setupMocks: func() {
				mockService.On("CreateProduct", mock.Anything, mock.Anything).Return(nil, pgutils.ErrDuplicatedKey).Once()
			},
			expectedCode: fiber.StatusConflict,
		},
	}
	runGeneralHandlerTests(t, tests, setupProductApp(mockService))
}

func TestProductHandler_updateProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPut,
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setupMocks: func() {
				mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(&dto.ProductOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "duplicate product",
			method:   fiber.MethodPut,
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setupMocks: func() {
				mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil, pgutils.ErrDuplicatedKey).Once()
			},
			expectedCode: fiber.StatusConflict,
		},
		{
			name:     "not found",
			method:   fiber.MethodPut,
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setupMocks: func() {
				mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupProductApp(mockService))
}

func TestProductHandler_deleteProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3]}`),
			setupMocks: func() {
				mockService.On("DeleteProducts", mock.Anything, []uint{1, 2, 3}).Return(nil).Once()
			},
			expectedCode: fiber.StatusNoContent,
		},
		{
			name:     "bad request",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4]}`),
			setupMocks: func() {
				mockService.On("DeleteProducts", mock.Anything, []uint{1, 2, 3, 4}).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5]}`),
			setupMocks: func() {
				mockService.On("DeleteProducts", mock.Anything, []uint{1, 2, 3, 4, 5}).Return(gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:         "invalid id",
			method:       fiber.MethodDelete,
			endpoint:     "/",
			body:         strings.NewReader(`{"ids": ["a", "b", "c", "d"]}`),
			setupMocks:   func() {},
			expectedCode: fiber.StatusBadRequest,
		},
	}
	runGeneralHandlerTests(t, tests, setupProductApp(mockService))
}
