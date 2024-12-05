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
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
)

func setupUserApp(mockService *mocks.UserServiceMock) *fiber.App {
	middleware.MidAccess = middleware.Auth(os.Getenv("ACCESS_TOKEN_PUBLIC"), &mocks.UserRepositoryMock{})

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
	NewUserHandler(app.Group("/user"), mockService)

	return app
}

func TestUserHandler_getUserPhoto(t *testing.T) {
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/user/1/photo",
			body:     nil,
			setupMocks: func() {
				mockService.On("GenerateUserPhotoURL", mock.Anything, uint(1)).Return("http://photo.url", nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "not found",
			method:   fiber.MethodGet,
			endpoint: "/user/200/photo",
			body:     nil,
			setupMocks: func() {
				mockService.On("GenerateUserPhotoURL", mock.Anything, uint(200)).Return("", gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
		{
			name:     "general error",
			method:   fiber.MethodGet,
			endpoint: "/user/500/photo",
			body:     nil,
			setupMocks: func() {
				mockService.On("GenerateUserPhotoURL", mock.Anything, uint(500)).Return("", errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

func TestUserHandler_getUsers(t *testing.T) {
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/user",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetUsers", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.UserOutputDTO]{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "general error",
			method:   fiber.MethodGet,
			endpoint: "/user",
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
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPost,
			endpoint: "/user",
			body:     strings.NewReader(`{"name":"user1","email":"example@email.com"}`),
			setupMocks: func() {
				mockService.On("CreateUser", mock.Anything, mock.Anything).Return(&dto.UserOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusCreated,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPost,
			endpoint: "/user",
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
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodGet,
			endpoint: "/user/1",
			body:     nil,
			setupMocks: func() {
				mockService.On("GetUserByID", mock.Anything, uint(1)).Return(&dto.UserOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "not found",
			method:   fiber.MethodGet,
			endpoint: "/user/200",
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
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPut,
			endpoint: "/user/1",
			body:     strings.NewReader(`{"name":"user1","email":"example@email.com"}`),
			setupMocks: func() {
				mockService.On("UpdateUser", mock.Anything, uint(1), mock.Anything).Return(&dto.UserOutputDTO{}, nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPut,
			endpoint: "/user/500",
			body:     strings.NewReader(`{"name":"user1"}`),
			setupMocks: func() {
				mockService.On("UpdateUser", mock.Anything, uint(500), mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodPut,
			endpoint: "/user/200",
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
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodDelete,
			endpoint: "/user",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4]}`),
			setupMocks: func() {
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3, 4}).Return(nil).Once()
			},
			expectedCode: fiber.StatusNoContent,
		},
		{
			name:     "bad request",
			method:   fiber.MethodDelete,
			endpoint: "/user",
			body:     strings.NewReader(`{"ids": [1, 2, 3, 4, 5]}`),
			setupMocks: func() {
				mockService.On("DeleteUsers", mock.Anything, []uint{1, 2, 3, 4, 5}).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodDelete,
			endpoint: "/user",
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
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodDelete,
			endpoint: "/user/pass?email=example1@email.com",
			body:     nil,
			setupMocks: func() {
				mockService.On("ResetUserPassword", mock.Anything, "example1@email.com").Return(nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:     "bad request",
			method:   fiber.MethodDelete,
			endpoint: "/user/pass?email=example2@email.com",
			body:     nil,
			setupMocks: func() {
				mockService.On("ResetUserPassword", mock.Anything, "example2@email.com").Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodDelete,
			endpoint: "/user/pass?email=example3@email.com",
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
	mockService := new(mocks.UserServiceMock)
	tests := []generalHandlerTest{
		{
			name:     "success",
			method:   fiber.MethodPut,
			endpoint: "/user/pass?email=example1@email.com",
			body:     strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD>"}`),
			setupMocks: func() {
				mockService.On("SetUserPassword", mock.Anything, "example1@email.com", mock.Anything).Return(nil).Once()
			},
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "different passwords",
			method:       fiber.MethodPut,
			endpoint:     "/user/pass?email=example1@email.com",
			body:         strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD2>"}`),
			setupMocks:   func() {},
			expectedCode: fiber.StatusBadRequest,
		},
		{
			name:     "bad request",
			method:   fiber.MethodPut,
			endpoint: "/user/pass?email=example2@email.com",
			body:     strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD>"}`),
			setupMocks: func() {
				mockService.On("SetUserPassword", mock.Anything, "example2@email.com", mock.Anything).Return(errors.New("error")).Once()
			},
			expectedCode: fiber.StatusInternalServerError,
		},
		{
			name:     "not found",
			method:   fiber.MethodPut,
			endpoint: "/user/pass?email=example3@email.com",
			body:     strings.NewReader(`{"password": "<PASSWORD>", "password_confirm": "<PASSWORD>"}`),
			setupMocks: func() {
				mockService.On("SetUserPassword", mock.Anything, "example3@email.com", mock.Anything).Return(gorm.ErrRecordNotFound).Once()
			},
			expectedCode: fiber.StatusNotFound,
		},
	}
	runGeneralHandlerTests(t, tests, setupUserApp(mockService))
}

//func TestUserHandler_setUserPhoto(t *testing.T) {
//	userHandler, userService := setup()
//
//	tests := []struct {
//		name    string
//		setup   func()
//		wantErr bool
//	}{
//		{
//			name: "success",
//			setup: func() {
//				userService.On("SetUserPhoto", mock.Anything, uint(1), mock.Anything).Return(nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("SetUserPhoto", mock.Anything, uint(1), mock.Anything).Return(errors.New("error"))
//			},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		tt.setup()
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			req := app.AcquireCtx(&fiber.Ctx{})
//			req.Locals(helper.LocalID, &filters.IDFilter{ID: 1})
//			req.Locals(helper.LocalDTO, &domain.File{})
//
//			err := userHandler.setUserPhoto(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusOK, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
