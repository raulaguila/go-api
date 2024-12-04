package handler

//
//import (
//	"context"
//	"errors"
//	"github.com/gofiber/fiber/v2"
//	"testing"
//
//	"github.com/stretchr/testify/mock"
//	"github.com/stretchr/testify/require"
//
//	"github.com/raulaguila/go-api/internal/api/middleware/datatransferobject"
//	"github.com/raulaguila/go-api/internal/pkg/domain"
//	"github.com/raulaguila/go-api/internal/pkg/dto"
//	"github.com/raulaguila/go-api/internal/pkg/filters"
//	"github.com/raulaguila/go-api/internal/pkg/myerrors"
//	"github.com/raulaguila/go-api/pkg/helper"
//)
//
//type MockUserService struct {
//	mock.Mock
//}
//
//func (m *MockUserService) GenerateUserPhotoURL(ctx context.Context, userID uint) (string, error) {
//	args := m.Called(ctx, userID)
//	return args.String(0), args.Error(1)
//}
//
//func (m *MockUserService) SetUserPhoto(ctx context.Context, userID uint, file *domain.File) error {
//	args := m.Called(ctx, userID, file)
//	return args.Error(0)
//}
//
//func (m *MockUserService) GetUsers(ctx context.Context, filter *filters.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error) {
//	args := m.Called(ctx, filter)
//	return args.Get(0).(*dto.ItemsOutputDTO[dto.UserOutputDTO]), args.Error(1)
//}
//
//func (m *MockUserService) CreateUser(ctx context.Context, userInput *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
//	args := m.Called(ctx, userInput)
//	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
//}
//
//func (m *MockUserService) GetUserByID(ctx context.Context, userID uint) (*dto.UserOutputDTO, error) {
//	args := m.Called(ctx, userID)
//	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
//}
//
//func (m *MockUserService) UpdateUser(ctx context.Context, userID uint, userInput *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
//	args := m.Called(ctx, userID, userInput)
//	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
//}
//
//func (m *MockUserService) DeleteUsers(ctx context.Context, userIDs []uint) error {
//	args := m.Called(ctx, userIDs)
//	return args.Error(0)
//}
//
//func (m *MockUserService) ResetUserPassword(ctx context.Context, email string) error {
//	args := m.Called(ctx, email)
//	return args.Error(0)
//}
//
//func (m *MockUserService) SetUserPassword(ctx context.Context, email string, passwordInput *dto.PasswordInputDTO) error {
//	args := m.Called(ctx, email, passwordInput)
//	return args.Error(0)
//}
//
//func setup() (userHandler *UserHandler, userService *MockUserService) {
//	userService = new(MockUserService)
//	userHandler = &UserHandler{
//		userService: userService,
//		handlerError: func(c *fiber.Ctx, err error) error {
//			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
//		},
//	}
//	return
//}
//
//func TestUserHandler_getUserPhoto(t *testing.T) {
//	userHandler, userService := setup()
//
//	tests := []struct {
//		name    string
//		setup   func()
//		wantURL string
//		wantErr bool
//	}{
//		{
//			name: "success",
//			setup: func() {
//				userService.On("GenerateUserPhotoURL", mock.Anything, uint(1)).Return("http://photo.url", nil)
//			},
//			wantURL: "http://photo.url",
//			wantErr: false,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("GenerateUserPhotoURL", mock.Anything, uint(1)).Return("", errors.New("error"))
//			},
//			wantURL: "",
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
//
//			err := userHandler.getUserPhoto(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, 302, req.Response().StatusCode())
//				require.Equal(t, tt.wantURL, req.Response().Header.Peek("Location"))
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
//
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
//
//func TestUserHandler_getUsers(t *testing.T) {
//	userHandler, userService := setup()
//
//	tests := []struct {
//		name     string
//		setup    func()
//		wantErr  bool
//		response *dto.ItemsOutputDTO[dto.UserOutputDTO]
//	}{
//		{
//			name: "success",
//			setup: func() {
//				userService.On("GetUsers", mock.Anything, mock.Anything).Return(&dto.ItemsOutputDTO[dto.UserOutputDTO]{}, nil)
//			},
//			wantErr:  false,
//			response: &dto.ItemsOutputDTO[dto.UserOutputDTO]{},
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("GetUsers", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
//			},
//			wantErr:  true,
//			response: nil,
//		},
//	}
//
//	for _, tt := range tests {
//		tt.setup()
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			req := app.AcquireCtx(&fiber.Ctx{})
//			req.Locals(helper.LocalFilter, &filters.UserFilter{})
//
//			err := userHandler.getUsers(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusOK, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
//
//func TestUserHandler_createUser(t *testing.T) {
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
//				userService.On("CreateUser", mock.Anything, mock.Anything).Return(&dto.UserOutputDTO{}, nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("error"))
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
//			req.Locals(helper.LocalDTO, &dto.UserInputDTO{})
//
//			err := userHandler.createUser(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusCreated, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
//
//func TestUserHandler_getUser(t *testing.T) {
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
//				userService.On("GetUserByID", mock.Anything, uint(1)).Return(&dto.UserOutputDTO{}, nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("GetUserByID", mock.Anything, uint(1)).Return(nil, errors.New("error"))
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
//
//			err := userHandler.getUser(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusOK, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
//
//func TestUserHandler_updateUser(t *testing.T) {
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
//				userService.On("UpdateUser", mock.Anything, uint(1), mock.Anything).Return(&dto.UserOutputDTO{}, nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("UpdateUser", mock.Anything, uint(1), mock.Anything).Return(nil, errors.New("error"))
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
//			req.Locals(helper.LocalDTO, &dto.UserInputDTO{})
//
//			err := userHandler.updateUser(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusOK, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
//
//func TestUserHandler_deleteUser(t *testing.T) {
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
//				userService.On("DeleteUsers", mock.Anything, mock.Anything).Return(nil)
//			},
//			wantErr: false,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("DeleteUsers", mock.Anything, mock.Anything).Return(errors.New("error"))
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
//			req.Locals(helper.LocalDTO, &dto.IDsInputDTO{IDs: []uint{1}})
//
//			err := userHandler.deleteUser(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusNoContent, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
//
//func TestUserHandler_resetUserPassword(t *testing.T) {
//	userHandler, userService := setup()
//
//	tests := []struct {
//		name    string
//		setup   func()
//		query   string
//		wantErr bool
//	}{
//		{
//			name: "success",
//			setup: func() {
//				userService.On("ResetUserPassword", mock.Anything, "email@example.com").Return(nil)
//			},
//			query:   "email%40example.com",
//			wantErr: false,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("ResetUserPassword", mock.Anything, "email@example.com").Return(errors.New("error"))
//			},
//			query:   "email%40example.com",
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		tt.setup()
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			req := app.AcquireCtx(&fiber.Ctx{})
//			req.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
//			req.Request().URI().SetQueryString(helper.ParamMail + "=" + tt.query)
//
//			err := userHandler.resetUserPassword(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusOK, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
//
//func TestUserHandler_setUserPassword(t *testing.T) {
//	userHandler, userService := setup()
//
//	tests := []struct {
//		name    string
//		setup   func()
//		query   string
//		dto     *dto.PasswordInputDTO
//		wantErr bool
//	}{
//		{
//			name: "success",
//			setup: func() {
//				userService.On("SetUserPassword", mock.Anything, "email@example.com", mock.Anything).Return(nil)
//			},
//			query: "email%40example.com",
//			dto: &dto.PasswordInputDTO{
//				Password:        &datatransferobject.Password{Value: "password"},
//				PasswordConfirm: &datatransferobject.Password{Value: "password"},
//			},
//			wantErr: false,
//		},
//		{
//			name: "password_nomatch",
//			setup: func() {
//				userService.On("SetUserPassword", mock.Anything, "email@example.com", mock.Anything).Return(myerrors.ErrPasswordsDoNotMatch)
//			},
//			query: "email%40example.com",
//			dto: &dto.PasswordInputDTO{
//				Password:        &datatransferobject.Password{Value: "password"},
//				PasswordConfirm: &datatransferobject.Password{Value: "different"},
//			},
//			wantErr: true,
//		},
//		{
//			name: "failure",
//			setup: func() {
//				userService.On("SetUserPassword", mock.Anything, "email@example.com", mock.Anything).Return(errors.New("error"))
//			},
//			query:   "email%40example.com",
//			dto:     &dto.PasswordInputDTO{Password: nil, PasswordConfirm: nil},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		tt.setup()
//		t.Run(tt.name, func(t *testing.T) {
//			app := fiber.New()
//			req := app.AcquireCtx(&fiber.Ctx{})
//			req.Set(fiber.HeaderContentType, fiber.MIMEApplicationJSON)
//			req.Request().URI().SetQueryString(helper.ParamMail + "=" + tt.query)
//			req.Locals(helper.LocalDTO, tt.dto)
//
//			err := userHandler.setUserPassword(req)
//			if tt.wantErr {
//				require.NotNil(t, err)
//			} else {
//				require.Equal(t, fiber.StatusOK, req.Response().StatusCode())
//			}
//		})
//		userService.AssertExpectations(t)
//	}
//}
