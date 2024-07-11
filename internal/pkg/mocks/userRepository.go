package mocks

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (s *UserRepositoryMock) CountUsers(ctx context.Context, f *filters.UserFilter) (int64, error) {
	ret := s.Called(ctx, f)
	return ret.Get(0).(int64), ret.Error(1)
}

func (s *UserRepositoryMock) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	ret := s.Called(ctx, userID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.User), ret.Error(1)
}

func (s *UserRepositoryMock) GetUsers(ctx context.Context, f *filters.UserFilter) (*[]domain.User, error) {
	ret := s.Called(ctx, f)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*[]domain.User), ret.Error(1)
}

func (s *UserRepositoryMock) CreateUser(ctx context.Context, data *dto.UserInputDTO) (*domain.User, error) {
	ret := s.Called(ctx, data)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.User), ret.Error(1)
}

func (s *UserRepositoryMock) UpdateUser(ctx context.Context, user *domain.User, data *dto.UserInputDTO) error {
	ret := s.Called(ctx, user, data)
	return ret.Error(0)
}

func (s *UserRepositoryMock) DeleteUsers(ctx context.Context, userIDs []uint) error {
	ret := s.Called(ctx, userIDs)
	return ret.Error(0)
}

func (s *UserRepositoryMock) GetUserByMail(ctx context.Context, mail string) (*domain.User, error) {
	ret := s.Called(ctx, mail)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.User), ret.Error(1)
}

func (s *UserRepositoryMock) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	ret := s.Called(ctx, token)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.User), ret.Error(1)
}

func (s *UserRepositoryMock) ResetUserPassword(ctx context.Context, user *domain.User) error {
	ret := s.Called(ctx, user)
	return ret.Error(0)
}

func (s *UserRepositoryMock) SetUserPassword(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	ret := s.Called(ctx, user, pass)
	return ret.Error(0)
}

func (s *UserRepositoryMock) SetUserPhoto(ctx context.Context, user *domain.User, file *domain.File) error {
	ret := s.Called(ctx, user, file)
	return ret.Error(0)
}

func (s *UserRepositoryMock) GenerateUserPhotoURL(ctx context.Context, user *domain.User) (string, error) {
	ret := s.Called(ctx, user)
	return ret.Get(0).(string), ret.Error(1)
}
