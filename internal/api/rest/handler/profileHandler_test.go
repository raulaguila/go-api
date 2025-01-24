package handler

import (
	"errors"
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
)

func setupProfileApp(mockService *mocks.ProfileServiceMock) *fiber.App {
	middleware.MidAccess = middleware.Auth(configs.AccessPrivateKey, &mocks.UserRepositoryMock{})

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
	NewProfileHandler(app.Group(""), mockService)

	return app
}

func TestProfileHandler_getProfiles(t *testing.T) {
	mockService := new(mocks.ProfileServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProfiles", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.ProfileOutputDTO]{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			method:   fiber.MethodGet,
			endpoint: "/",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProfiles", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}
	runGeneralHandlerTests(t, tests, setupProfileApp(mockService))
}

func TestProfileHandler_createProfile(t *testing.T) {
	mockService := new(mocks.ProfileServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"name":"admin"}`),
			setupMocks: func() {
				mockService.On("CreateProfile", mock.Anything, mock.Anything).Return(&dto.ProfileOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"name":"admin"}`),
			setupMocks: func() {
				mockService.On("CreateProfile", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}
	runGeneralHandlerTests(t, tests, setupProfileApp(mockService))
}

func TestProfileHandler_getProfile(t *testing.T) {
	mockService := new(mocks.ProfileServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/1",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProfileByID", mock.Anything, uint(1)).Return(&dto.ProfileOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "not found",
			method:   fiber.MethodGet,
			endpoint: "/200",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetProfileByID", mock.Anything, uint(200)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupProfileApp(mockService))
}

func TestProfileHandler_updateProfile(t *testing.T) {
	mockService := new(mocks.ProfileServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPut,
			endpoint: "/1",
			body:     strings.NewReader(`{"name":"user1","email":"example@email.com"}`),
			setupMocks: func() {
				mockService.On("UpdateProfile", mock.Anything, uint(1), mock.Anything).Return(&dto.ProfileOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPut,
			endpoint: "/500",
			body:     strings.NewReader(`{"name":"user1"}`),
			setupMocks: func() {
				mockService.On("UpdateProfile", mock.Anything, uint(500), mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodPut,
			endpoint: "/200",
			body:     strings.NewReader(`{"name":"user1"}`),
			setupMocks: func() {
				mockService.On("UpdateProfile", mock.Anything, uint(200), mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupProfileApp(mockService))
}

func TestProfileHandler_deleteProfile(t *testing.T) {
	mockService := new(mocks.ProfileServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4]}`),
			setupMocks: func() {
				mockService.On("DeleteProfiles", mock.Anything, []uint{1, 2, 3, 4}).Return(nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5]}`),
			setupMocks: func() {
				mockService.On("DeleteProfiles", mock.Anything, []uint{1, 2, 3, 4, 5}).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5, 6]}`),
			setupMocks: func() {
				mockService.On("DeleteProfiles", mock.Anything, []uint{1, 2, 3, 4, 5, 6}).Return(gorm.ErrRecordNotFound).Once()
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
	runGeneralHandlerTests(t, tests, setupProfileApp(mockService))
}
