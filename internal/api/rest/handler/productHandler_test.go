package handler

import (
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgerror"
	"github.com/raulaguila/packhub"
)

func setupProductApp(mockService *_mocks.ProductServiceMock) *fiber.App {
	middleware.MidAccess = middleware.Auth(configs.AccessPrivateKey, &_mocks.UserRepositoryMock{})

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
	tests := []struct {
		name, endpoint string
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			setup: func() {
				mockService.On("GetProducts", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.ProductOutputDTO]{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			endpoint: "/",
			setup: func() {
				mockService.On("GetProducts", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}

	app := setupProductApp(mockService)
	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			tt.setup()
			req := httptest.NewRequest(fiber.MethodGet, tt.endpoint, nil)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := app.Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_getProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []struct {
		name, endpoint string
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/1",
			setup: func() {
				mockService.On("GetProductByID", mock.Anything, uint(1)).Return(&dto.ProductOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			endpoint: "/200",
			setup: func() {
				mockService.On("GetProductByID", mock.Anything, uint(200)).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			endpoint: "/500",
			setup: func() {
				mockService.On("GetProductByID", mock.Anything, uint(500)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}

	app := setupProductApp(mockService)
	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			tt.setup()
			req := httptest.NewRequest(fiber.MethodGet, tt.endpoint, nil)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := app.Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_createProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() {
				mockService.On("CreateProduct", mock.Anything, mock.Anything).Return(&dto.ProductOutputDTO{
					ID:   packhub.Pointer(uint(1)),
					Name: packhub.Pointer("Product 01"),
				}, nil).Once()
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "duplicate product",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() {
				mockService.On("CreateProduct", mock.Anything, mock.Anything).Return(nil, pgerror.ErrDuplicatedKey).Once()
			},
			expectedCode: fiber.StatusConflict,
		},
	}

	app := setupProductApp(mockService)
	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			tt.setup()
			req := httptest.NewRequest(fiber.MethodPost, tt.endpoint, tt.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := app.Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_updateProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() {
				mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(&dto.ProductOutputDTO{
					ID:   packhub.Pointer(uint(1)),
					Name: packhub.Pointer("Product 01"),
				}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "duplicate product",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() {
				mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil, pgerror.ErrDuplicatedKey).Once()
			},
			expectedCode: fiber.StatusConflict,
		},
		{
			name:     "not found",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "Product 01"}`),
			setup: func() {
				mockService.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}

	app := setupProductApp(mockService)
	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			tt.setup()
			req := httptest.NewRequest(fiber.MethodPut, tt.endpoint, tt.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := app.Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}

func TestProductHandler_deleteProduct(t *testing.T) {
	mockService := new(_mocks.ProductServiceMock)
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3]}`),
			setup: func() {
				mockService.On("DeleteProducts", mock.Anything, []uint{1, 2, 3}).Return(nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4]}`),
			setup: func() {
				mockService.On("DeleteProducts", mock.Anything, []uint{1, 2, 3, 4}).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5]}`),
			setup: func() {
				mockService.On("DeleteProducts", mock.Anything, []uint{1, 2, 3, 4, 5}).Return(gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:         "invalid id",
			endpoint:     "/",
			body:         strings.NewReader(`{"ids": ["a", "b", "c", "d"]}`),
			setup:        func() {},
			expectedCode: fiber.StatusBadRequest,
		},
	}

	app := setupProductApp(mockService)
	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			tt.setup()
			req := httptest.NewRequest(fiber.MethodDelete, tt.endpoint, tt.body)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := app.Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}
