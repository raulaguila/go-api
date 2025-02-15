package handler

import (
	"errors"
	"fmt"
	"golang.org/x/text/language"
	"gorm.io/gorm"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgerror"
	"github.com/raulaguila/packhub"
)

func setupUserApp(mockService *_mocks.UserServiceMock) *fiber.App {
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
	NewUserHandler(app.Group(""), mockService)

	return app
}

func TestUserHandler_getUsers(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []struct {
		name, endpoint string
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			setup: func() {
				mockService.On("GetUsers", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.UserOutputDTO]{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			endpoint: "/",
			setup: func() {
				mockService.On("GetUsers", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}

	app := setupUserApp(mockService)
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

func TestUserHandler_getUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []struct {
		name, endpoint string
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/1",
			setup: func() {
				mockService.On("GetUserByID", mock.Anything, uint(1)).Return(&dto.UserOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			endpoint: "/200",
			setup: func() {
				mockService.On("GetUserByID", mock.Anything, uint(200)).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			endpoint: "/500",
			setup: func() {
				mockService.On("GetUserByID", mock.Anything, uint(500)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}

	app := setupUserApp(mockService)
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

func TestUserHandler_createUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "User 01"}`),
			setup: func() {
				mockService.On("CreateUser", mock.Anything, mock.Anything).Return(&dto.UserOutputDTO{
					ID:   packhub.Pointer(uint(1)),
					Name: packhub.Pointer("User 01"),
				}, nil).Once()
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "duplicate user",
			endpoint: "/",
			body:     strings.NewReader(`{"name": "User 01"}`),
			setup: func() {
				mockService.On("CreateUser", mock.Anything, mock.Anything).Return(nil, pgerror.ErrDuplicatedKey).Once()
			},
			expectedCode: fiber.StatusConflict,
		},
	}

	app := setupUserApp(mockService)
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

func TestUserHandler_updateUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func()
		expectedCode   int
	}{
		{
			name:     "success",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "User 01"}`),
			setup: func() {
				mockService.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(&dto.UserOutputDTO{
					ID:   packhub.Pointer(uint(1)),
					Name: packhub.Pointer("User 01"),
				}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "duplicate user",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "User 01"}`),
			setup: func() {
				mockService.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, pgerror.ErrDuplicatedKey).Once()
			},
			expectedCode: fiber.StatusConflict,
		},
		{
			name:     "not found",
			endpoint: "/1",
			body:     strings.NewReader(`{"name": "User 01"}`),
			setup: func() {
				mockService.On("UpdateUser", mock.Anything, mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}

	app := setupUserApp(mockService)
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

func TestUserHandler_deleteUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
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
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3}).Return(nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4]}`),
			setup: func() {
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3, 4}).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5]}`),
			setup: func() {
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3, 4, 5}).Return(gorm.ErrRecordNotFound).Once()
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

	app := setupUserApp(mockService)
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
