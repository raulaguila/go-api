package _mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
)

func NewProfileServiceMock() domain.ProfileService {
	return new(ProfileServiceMock)
}

type ProfileServiceMock struct {
	mock.Mock
}

func (m *ProfileServiceMock) GenerateProfileOutputDTO(profile *domain.Profile) *dto.ProfileOutputDTO {
	args := m.Called(profile)
	return args.Get(0).(*dto.ProfileOutputDTO)
}

func (m *ProfileServiceMock) GetProfileByID(ctx context.Context, id uint) (*dto.ProfileOutputDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ProfileOutputDTO), args.Error(1)
}

func (m *ProfileServiceMock) GetProfiles(ctx context.Context, filter *dto.ProfileFilter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ItemsOutputDTO[dto.ProfileOutputDTO]), args.Error(1)
}

func (m *ProfileServiceMock) CreateProfile(ctx context.Context, input *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ProfileOutputDTO), args.Error(1)
}

func (m *ProfileServiceMock) UpdateProfile(ctx context.Context, id uint, input *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	args := m.Called(ctx, id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*dto.ProfileOutputDTO), args.Error(1)
}

func (m *ProfileServiceMock) DeleteProfiles(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}
