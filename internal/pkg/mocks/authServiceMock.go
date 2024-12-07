package mocks

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) Login(ctx context.Context, input *dto.AuthInputDTO, ip string) (*dto.AuthOutputDTO, error) {
	args := m.Called(ctx, input, ip)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.AuthOutputDTO), args.Error(1)
}

func (m *AuthServiceMock) Refresh(user *domain.User, ip string) *dto.AuthOutputDTO {
	args := m.Called(user, ip)
	return args.Get(0).(*dto.AuthOutputDTO)
}

func (m *AuthServiceMock) Me(user *domain.User) *dto.UserOutputDTO {
	args := m.Called(user)
	return args.Get(0).(*dto.UserOutputDTO)
}
