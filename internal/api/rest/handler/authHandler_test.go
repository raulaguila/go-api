package handler

import (
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"github.com/raulaguila/go-api/pkg/utils"
)

func setupAuthApp(nockService *mocks.AuthServiceMock) *fiber.App {
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
	nockService := new(mocks.AuthServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "valid login",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setupMocks: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(&dto.AuthOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "invalid body",
			method:       fiber.MethodPost,
			endpoint:     "/",
			body:         strings.NewReader(`{invalidJson}`),
			setupMocks:   func() {},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name:     "not found",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setupMocks: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:     "disabled user",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setupMocks: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(nil, utils.ErrDisabledUser).Once()
			},
			expectedCode: fiber.StatusUnauthorized,
		},
		{
			name:     "incorrect credentials",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setupMocks: func() {
				nockService.On("Login", mock.Anything, mock.Anything).Return(nil, utils.ErrInvalidCredentials).Once()
			},
			expectedCode: fiber.StatusUnauthorized,
		},
	}
	runGeneralHandlerTests(t, tests, setupAuthApp(nockService))
}

func TestAuthHandler_me(t *testing.T) {
	middleware.MidAccess = middleware.Auth(configs.AccessPrivateKey, &mocks.UserRepositoryMock{})

	nockService := new(mocks.AuthServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "valid request",
			method:   fiber.MethodGet,
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setupMocks: func() {
				nockService.On("Me", mock.Anything).Return(&dto.UserOutputDTO{}).Once()
			},
			expectedCode: fiber.StatusOK,
		},
	}
	runGeneralHandlerTests(t, tests, setupAuthApp(nockService))
}

func TestAuthHandler_refresh(t *testing.T) {
	middleware.MidRefresh = middleware.Auth(configs.RefreshPrivateKey, &mocks.UserRepositoryMock{})

	nockService := new(mocks.AuthServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "valid request",
			method:   fiber.MethodPut,
			endpoint: "/",
			body:     strings.NewReader(`{"username":"admin@admin.com","password":"12345678"}`),
			setupMocks: func() {
				nockService.On("Refresh", mock.Anything).Return(&dto.AuthOutputDTO{}).Once()
			},
			expectedCode: fiber.StatusOK,
		},
	}
	runGeneralHandlerTests(t, tests, setupAuthApp(nockService))
}
