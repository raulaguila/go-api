package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/contrib/fiberi18n/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/raulaguila/go-api/configs"
	_ "github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"golang.org/x/text/language"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

type MockProfileService struct {
	mock.Mock
}

func (m *MockProfileService) GenerateProfileOutputDTO(profile *domain.Profile) *dto.ProfileOutputDTO {
	args := m.Called(profile)
	return args.Get(0).(*dto.ProfileOutputDTO)
}

func (m *MockProfileService) GetProfileByID(ctx context.Context, id uint) (*dto.ProfileOutputDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ProfileOutputDTO), args.Error(1)
}

func (m *MockProfileService) GetProfiles(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ItemsOutputDTO[dto.ProfileOutputDTO]), args.Error(1)
}

func (m *MockProfileService) CreateProfile(ctx context.Context, input *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ProfileOutputDTO), args.Error(1)
}

func (m *MockProfileService) UpdateProfile(ctx context.Context, id uint, input *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ProfileOutputDTO), args.Error(1)
}

func (m *MockProfileService) DeleteProfiles(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func TestGetProfiles(t *testing.T) {
	mockService := new(MockProfileService)
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
	NewProfileHandler(app.Group("/profile"), mockService)

	tests := []struct {
		name       string
		setupMocks func()
		expectCode int
	}{
		{
			name: "Success",
			setupMocks: func() {
				mockService.On("GetProfiles", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.ProfileOutputDTO]{}, nil).Once()
			},
			expectCode: fiber.StatusOK,
		},
		{
			name: "Error",
			setupMocks: func() {
				mockService.On("GetProfiles", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectCode: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			req := httptest.NewRequest(fiber.MethodGet, "/profile", nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectCode, resp.StatusCode)
		})
	}
}

func TestGetProfile(t *testing.T) {
	mockService := new(MockProfileService)
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
	NewProfileHandler(app.Group("/profile"), mockService)

	tests := []struct {
		name       string
		id         uint
		setupMocks func()
		expectCode int
	}{
		{
			name: "Success",
			id:   1,
			setupMocks: func() {
				mockService.On("GetProfileByID", mock.Anything, mock.Anything).Return(&dto.ProfileOutputDTO{}, nil).Once()
			},
			expectCode: fiber.StatusOK,
		},
		{
			name: "Error",
			id:   200,
			setupMocks: func() {
				mockService.On("GetProfileByID", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectCode: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			req := httptest.NewRequest(fiber.MethodGet, fmt.Sprintf("/profile/%v", tt.id), nil)
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectCode, resp.StatusCode)
		})
	}
}

func TestCreateProfile(t *testing.T) {
	mockService := new(MockProfileService)
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
	NewProfileHandler(app.Group("/profile"), mockService)

	tests := []struct {
		name       string
		body       string
		setupMocks func()
		expectCode int
	}{
		{
			name: "Success",
			body: `{"name":"Admin"}`,
			setupMocks: func() {
				mockService.On("CreateProfile", mock.Anything, mock.Anything).Return(&dto.ProfileOutputDTO{}, nil).Once()
			},
			expectCode: fiber.StatusCreated,
		},
		{
			name: "Error",
			body: `{"name":"Admin"}`,
			setupMocks: func() {
				mockService.On("CreateProfile", mock.Anything, mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectCode: fiber.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			req := httptest.NewRequest(fiber.MethodPost, "/profile", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectCode, resp.StatusCode)
		})
	}
}

//func TestUpdateProfile(t *testing.T) {
//	mockService := new(MockProfileService)
//	handler := &ProfileHandler{
//		profileService: mockService,
//		handlerError: func(c *fiber.Ctx, err error) error {
//			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
//		},
//	}
//
//	tests := []struct {
//		name       string
//		setupMocks func()
//		expectCode int
//	}{
//		{
//			name: "Success",
//			setupMocks: func() {
//				mockService.On("UpdateProfile", mock.Anything, mock.Anything, mock.Anything).Return(&dto.ProfileOutputDTO{}, nil)
//			},
//			expectCode: fiber.StatusOK,
//		},
//		{
//			name: "Error",
//			setupMocks: func() {
//				mockService.On("UpdateProfile", mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("error"))
//			},
//			expectCode: fiber.StatusInternalServerError,
//		},
//	}
//
//	app := fiber.New()
//	app.Put("/:id", func(c *fiber.Ctx) error {
//		return handler.updateProfile(c)
//	})
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMocks()
//			req := fiber.Request{}
//			req.Header.SetMethod(fiber.MethodPut)
//			resp, err := app.Test(req)
//			assert.NoError(t, err)
//			assert.Equal(t, tt.expectCode, resp.StatusCode)
//			mockService.AssertExpectations(t)
//		})
//	}
//}
//
//func TestDeleteProfiles(t *testing.T) {
//	mockService := new(MockProfileService)
//	handler := &ProfileHandler{
//		profileService: mockService,
//		handlerError: func(c *fiber.Ctx, err error) error {
//			return c.Status(fiber.StatusInternalServerError).JSON(err.Error())
//		},
//	}
//
//	tests := []struct {
//		name       string
//		setupMocks func()
//		expectCode int
//	}{
//		{
//			name: "Success",
//			setupMocks: func() {
//				mockService.On("DeleteProfiles", mock.Anything, mock.Anything).Return(nil)
//			},
//			expectCode: fiber.StatusNoContent,
//		},
//		{
//			name: "Error",
//			setupMocks: func() {
//				mockService.On("DeleteProfiles", mock.Anything, mock.Anything).Return(errors.New("error"))
//			},
//			expectCode: fiber.StatusInternalServerError,
//		},
//	}
//
//	app := fiber.New()
//	app.Delete("/", func(c *fiber.Ctx) error {
//		return handler.deleteProfiles(c)
//	})
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMocks()
//			req := fiber.Request{}
//			req.Header.SetMethod(fiber.MethodDelete)
//			resp, err := app.Test(req)
//			assert.NoError(t, err)
//			assert.Equal(t, tt.expectCode, resp.StatusCode)
//			mockService.AssertExpectations(t)
//		})
//	}
//}
