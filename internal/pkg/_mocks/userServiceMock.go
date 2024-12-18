package _mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
)

type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) GenerateUserOutputDTO(profile *domain.User) *dto.UserOutputDTO {
	args := m.Called(profile)
	return args.Get(0).(*dto.UserOutputDTO)
}

func (m *UserServiceMock) GenerateUserPhotoURL(ctx context.Context, userID uint) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

func (m *UserServiceMock) SetUserPhoto(ctx context.Context, userID uint, file *domain.File) error {
	args := m.Called(ctx, userID, file)
	return args.Error(0)
}

func (m *UserServiceMock) GetUsers(ctx context.Context, filter *filters.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.ItemsOutputDTO[dto.UserOutputDTO]), args.Error(1)
}

func (m *UserServiceMock) CreateUser(ctx context.Context, userInput *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	args := m.Called(ctx, userInput)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

func (m *UserServiceMock) GetUserByID(ctx context.Context, userID uint) (*dto.UserOutputDTO, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

func (m *UserServiceMock) UpdateUser(ctx context.Context, userID uint, userInput *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	args := m.Called(ctx, userID, userInput)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserOutputDTO), args.Error(1)
}

func (m *UserServiceMock) DeleteUsers(ctx context.Context, userIDs []uint) error {
	args := m.Called(ctx, userIDs)
	return args.Error(0)
}

func (m *UserServiceMock) ResetUserPassword(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *UserServiceMock) SetUserPassword(ctx context.Context, email string, passwordInput *dto.PasswordInputDTO) error {
	args := m.Called(ctx, email, passwordInput)
	return args.Error(0)
}
