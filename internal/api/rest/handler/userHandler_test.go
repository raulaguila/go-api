package handler

import (
	"errors"
	"github.com/raulaguila/go-api/internal/api/rest/middleware"
	"strings"
	"testing"

	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/mock"
	"golang.org/x/text/language"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/dto"
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
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetUsers", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.UserOutputDTO]{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			method:   fiber.MethodGet,
			endpoint: "/",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetUsers", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

func TestUserHandler_createUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"name":"user1","email":"example@email.com"}`),
			setupMocks: func() {
				mockService.On("CreateUser", mock.Anything, mock.Anything).Return(&dto.UserOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPost,
			endpoint: "/",
			body:     strings.NewReader(`{"name":"user1"}`),
			setupMocks: func() {
				mockService.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

func TestUserHandler_getUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/1",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetUserByID", mock.Anything, uint(1)).Return(&dto.UserOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "not found",
			method:   fiber.MethodGet,
			endpoint: "/200",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetUserByID", mock.Anything, uint(200)).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

func TestUserHandler_updateUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPut,
			endpoint: "/1",
			body:     strings.NewReader(`{"name":"user1","email":"example@email.com"}`),
			setupMocks: func() {
				mockService.On("UpdateUser", mock.Anything, uint(1), mock.Anything).Return(&dto.UserOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPut,
			endpoint: "/500",
			body:     strings.NewReader(`{"name":"user1"}`),
			setupMocks: func() {
				mockService.On("UpdateUser", mock.Anything, uint(500), mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodPut,
			endpoint: "/200",
			body:     strings.NewReader(`{"name":"user1"}`),
			setupMocks: func() {
				mockService.On("UpdateUser", mock.Anything, uint(200), mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

func TestUserHandler_deleteUser(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4]}`),
			setupMocks: func() {
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3, 4}).Return(nil).Once()
			},
			expectedCode: fiber.StatusNoContent,
		},
		{
			name:     "bad request",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5]}`),
			setupMocks: func() {
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3, 4, 5}).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodDelete,
			endpoint: "/",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5, 6]}`),
			setupMocks: func() {
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3, 4, 5, 6}).Return(gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

func TestUserHandler_resetUserPassword(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodDelete,
			endpoint: "/pass?email=example1@email.com",
			body:     nil,
			setupMocks: func() {
				mockService.On("ResetUserPassword", mock.Anything, "example1@email.com").Return(nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			method:   fiber.MethodDelete,
			endpoint: "/pass?email=example2@email.com",
			body:     nil,
			setupMocks: func() {
				mockService.On("ResetUserPassword", mock.Anything, "example2@email.com").Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodDelete,
			endpoint: "/pass?email=example3@email.com",
			body:     nil,
			setupMocks: func() {
				mockService.On("ResetUserPassword", mock.Anything, "example3@email.com").Return(gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

func TestUserHandler_setUserPassword(t *testing.T) {
	mockService := new(_mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPut,
			endpoint: "/pass?email=example1@email.com",
			body:     strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD>"}`),
			setupMocks: func() {
				mockService.On("SetUserPassword", mock.Anything, "example1@email.com", mock.Anything).Return(nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "different passwords",
			method:       fiber.MethodPut,
			endpoint:     "/pass?email=example1@email.com",
			body:         strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD2>"}`),
			setupMocks:   func() {},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPut,
			endpoint: "/pass?email=example2@email.com",
			body:     strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD>"}`),
			setupMocks: func() {
				mockService.On("SetUserPassword", mock.Anything, "example2@email.com", mock.Anything).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodPut,
			endpoint: "/pass?email=example3@email.com",
			body:     strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD>"}`),
			setupMocks: func() {
				mockService.On("SetUserPassword", mock.Anything, "example3@email.com", mock.Anything).Return(gorm.ErrRecordNotFound).Once()
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
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}
