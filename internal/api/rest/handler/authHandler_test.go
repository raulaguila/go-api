package handler

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/utils"
)

func setupAuthApp(nockService *_mocks.AuthServiceMock) *fiber.App {
	app := fiber.New()
	app.Use(fiberi18n.New(&fiberi18n.Config{
		Next:            func(c *fiber.Ctx) bool { return false },
		RootPath:        "./locales",
		AcceptLanguages: []language.Tag{language.AmericanEnglish, language.BrazilianPortuguese},
		DefaultLanguage: language.AmericanEnglish,
		Loader:          &fiberi18n.EmbedLoader{FS: configs.Locales},
	}))
	NewAuthHandler(app.Group(""), nockService)

	return app
}

func TestAuthHandler_login(t *testing.T) {
	nockService := new(_mocks.AuthServiceMock)
	tests := []struct {
		name, endpoint string
		body           io.Reader
		setup          func()
		expectedCode   int
	}{
		{
			name:     "valid login",
			endpoint: "/",
			body:     strings.NewReader(`{"login": "admin@admin.com","password": "12345678"}`),
			setup: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(&dto.AuthOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "invalid body",
			endpoint:     "/",
			body:         strings.NewReader(`{invalidJson}`),
			setup:        func() {},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name:     "not found",
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setup: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:     "disabled user",
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setup: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(nil, utils.ErrDisabledUser).Once()
			},
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			name:     "incorrect credentials",
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setup: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(nil, utils.ErrInvalidCredentials).Once()
			},
			expectedCode: fiber.StatusUnauthorized,
		},
	}

	app := setupAuthApp(nockService)
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

func TestAuthHandler_me(t *testing.T) {
	middleware.MidAccess = middleware.Auth(configs.AccessPrivateKey, &_mocks.UserRepositoryMock{})

	nockService := new(_mocks.AuthServiceMock)
	tests := []struct {
		name, endpoint string
		setup          func()
		expectedCode   int
	}{
		{
			name:     "valid request",
			endpoint: "/",
			setup: func() {
				nockService.On("Me", mock.Anything).Return(&dto.UserOutputDTO{}).Once()
			},
			expectedCode: fiber.StatusOK,
		},
	}

	app := setupAuthApp(nockService)
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

func TestAuthHandler_refresh(t *testing.T) {
	middleware.MidRefresh = middleware.Auth(configs.RefreshPrivateKey, &_mocks.UserRepositoryMock{})

	nockService := new(_mocks.AuthServiceMock)
	tests := []struct {
		name, endpoint string
		setup          func()
		expectedCode   int
	}{
		{
			name:     "valid request",
			endpoint: "/",
			setup: func() {
				nockService.On("Refresh", mock.Anything).Return(&dto.AuthOutputDTO{}).Once()
			},
			expectedCode: fiber.StatusOK,
		},
	}

	app := setupAuthApp(nockService)
	for _, tt := range tests {
		t.Run(tt.name, func(test *testing.T) {
			tt.setup()
			req := httptest.NewRequest(fiber.MethodPut, tt.endpoint, nil)
			req.Header.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
			req.Header.Set("X-Skip-Auth", "true")

			resp, err := app.Test(req)
			require.NoError(test, err, fmt.Sprintf("Error on test '%v'", tt.name))
			require.Equal(test, tt.expectedCode, resp.StatusCode, fmt.Sprintf("Wrong status code on test '%v'", tt.name))
		})
	}
}
