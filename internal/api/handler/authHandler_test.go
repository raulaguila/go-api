package handler

//
//import (
//	"bytes"
//	"context"
//	"github.com/gofiber/fiber/v2"
//	"github.com/raulaguila/go-api/pkg/helper"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/raulaguila/go-api/internal/pkg/domain"
//	"github.com/raulaguila/go-api/internal/pkg/dto"
//	"github.com/raulaguila/go-api/internal/pkg/myerrors"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//	"gorm.io/gorm"
//)
//
//type MockAuthService struct {
//	mock.Mock
//}
//
//func (m *MockAuthService) Login(ctx context.Context, input *dto.AuthInputDTO, ip string) (*dto.AuthOutputDTO, error) {
//	args := m.Called(ctx, input, ip)
//	return args.Get(0).(*dto.AuthOutputDTO), args.Error(1)
//}
//
//func (m *MockAuthService) Refresh(user *domain.User, ip string) *dto.AuthOutputDTO {
//	args := m.Called(user, ip)
//	return args.Get(0).(*dto.AuthOutputDTO)
//}
//
//func (m *MockAuthService) Me(user *domain.User) *dto.UserOutputDTO {
//	args := m.Called(user)
//	return args.Get(0).(*dto.UserOutputDTO)
//}
//
//func TestAuthHandler_login(t *testing.T) {
//	app := fiber.New()
//	mockAuthService := new(MockAuthService)
//	handler := &AuthHandler{
//		authService: mockAuthService,
//		handlerError: newErrorHandler(map[string]map[error][]any{
//			"*": {
//				myerrors.ErrDisabledUser:       []any{fiber.StatusUnauthorized, "disabledUser"},
//				myerrors.ErrInvalidCredentials: []any{fiber.StatusUnauthorized, "incorrectCredentials"},
//				gorm.ErrRecordNotFound:         []any{fiber.StatusNotFound, "userNotFound"},
//			},
//		}),
//	}
//
//	app.Post("/", handler.login)
//
//	tests := []struct {
//		name         string
//		body         string
//		wantStatus   int
//		authResponse *dto.AuthOutputDTO
//		mockError    error
//	}{
//		{
//			name:       "valid login",
//			body:       `{"username":"testuser","password":"testpass"}`,
//			wantStatus: fiber.StatusOK,
//			authResponse: &dto.AuthOutputDTO{
//				Token: "validToken",
//			},
//			mockError: nil,
//		},
//		{
//			name:       "invalid body",
//			body:       `{invalidJson}`,
//			wantStatus: fiber.StatusBadRequest,
//		},
//		{
//			name:         "login error",
//			body:         `{"username":"testuser","password":"wrongpass"}`,
//			wantStatus:   fiber.StatusUnauthorized,
//			authResponse: nil,
//			mockError:    myerrors.ErrInvalidCredentials,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(tt.body)))
//			req.Header.Set("Content-Type", "application/json")
//
//			if tt.authResponse != nil || tt.mockError != nil {
//				mockAuthService.On("Login", mock.Anything, mock.Anything, mock.Anything).Return(tt.authResponse, tt.mockError)
//			}
//
//			resp, err := app.Test(req)
//
//			assert.NoError(t, err)
//			assert.Equal(t, tt.wantStatus, resp.StatusCode)
//		})
//	}
//}
//
//func TestAuthHandler_me(t *testing.T) {
//	app := fiber.New()
//	mockAuthService := new(MockAuthService)
//	handler := &AuthHandler{authService: mockAuthService}
//
//	app.Get("/", func(c *fiber.Ctx) error {
//		c.Locals(helper.LocalUser, &domain.User{ID: 1})
//		return handler.me(c)
//	})
//
//	tests := []struct {
//		name       string
//		user       *domain.User
//		wantStatus int
//		meResponse *dto.UserOutputDTO
//	}{
//		{
//			name:       "valid me request",
//			user:       &domain.User{ID: 1},
//			wantStatus: fiber.StatusOK,
//			meResponse: &dto.UserOutputDTO{ID: 1},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mockAuthService.On("Me", tt.user).Return(tt.meResponse)
//
//			req := httptest.NewRequest("GET", "/", nil)
//
//			resp, err := app.Test(req)
//			assert.NoError(t, err)
//			assert.Equal(t, tt.wantStatus, resp.StatusCode)
//		})
//	}
//}
//
//func TestAuthHandler_refresh(t *testing.T) {
//	app := fiber.New()
//	mockAuthService := new(MockAuthService)
//	handler := &AuthHandler{authService: mockAuthService}
//
//	app.Put("/", func(c *fiber.Ctx) error {
//		c.Locals(helper.LocalUser, &domain.User{ID: 1})
//		return handler.refresh(c)
//	})
//
//	tests := []struct {
//		name          string
//		user          *domain.User
//		wantStatus    int
//		refreshResult *dto.AuthOutputDTO
//	}{
//		{
//			name:          "valid refresh request",
//			user:          &domain.User{ID: 1},
//			wantStatus:    fiber.StatusOK,
//			refreshResult: &dto.AuthOutputDTO{Token: "newToken"},
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			mockAuthService.On("Refresh", tt.user, mock.Anything).Return(tt.refreshResult)
//
//			req := httptest.NewRequest("PUT", "/", nil)
//
//			resp, err := app.Test(req)
//			assert.NoError(t, err)
//			assert.Equal(t, tt.wantStatus, resp.StatusCode)
//		})
//	}
//}
