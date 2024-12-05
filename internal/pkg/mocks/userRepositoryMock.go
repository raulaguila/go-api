package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *UserRepositoryMock) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) GetUsers(ctx context.Context, userFilter *filters.UserFilter) (*[]domain.User, error) {
	args := m.Called(ctx, userFilter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]domain.User), args.Error(1)
}

func (m *UserRepositoryMock) CountUsers(ctx context.Context, userFilter *filters.UserFilter) (int64, error) {
	args := m.Called(ctx, userFilter)
	return args.Get(0).(int64), args.Error(1)
}

func (m *UserRepositoryMock) CreateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) DeleteUsers(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func (m *UserRepositoryMock) ResetUserPassword(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) SetUserPassword(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	args := m.Called(ctx, user, pass)
	return args.Error(0)
}

func (m *UserRepositoryMock) GetUserByMail(ctx context.Context, mail string) (*domain.User, error) {
	args := m.Called(ctx, mail)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *UserRepositoryMock) SetUserPhoto(ctx context.Context, user *domain.User, p *domain.File) error {
	args := m.Called(ctx, user, p)
	return args.Error(0)
}

func (m *UserRepositoryMock) GenerateUserPhotoURL(ctx context.Context, user *domain.User) (string, error) {
	args := m.Called(ctx, user)
	return args.String(0), args.Error(1)
}
