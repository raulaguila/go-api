package _mocks

import (
	"context"
	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
)

func NewUserRepositoryMock() domain.UserRepository {
	return new(UserRepositoryMock)
}

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetUser(ctx context.Context, user *domain.User) error {
	user.Name = "John Doe"
	user.Email = "johndoe@example.com"
	user.Auth = &domain.Auth{
		Status: false,
		Profile: &domain.Profile{
			Base: domain.Base{ID: uint(1)},
			Name: "ADMIN",
		},
	}
	args := m.Called(ctx, user)
	return args.Error(0)
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
	user.ID = 1
	user.Name = "John Doe"
	user.Email = "johndoe@example.com"
	user.Auth = &domain.Auth{Status: true, Profile: &domain.Profile{Base: domain.Base{ID: 1}, Name: "ADMIN"}}

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
