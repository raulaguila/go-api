package mocks

import (
	"context"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/stretchr/testify/mock"
)

type ProfileRepositoryMock struct {
	mock.Mock
}

func (s *ProfileRepositoryMock) CountProfiles(ctx context.Context, f *filter.Filter) (int64, error) {
	ret := s.Called(ctx, f)
	return ret.Get(0).(int64), ret.Error(1)
}

func (s *ProfileRepositoryMock) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	ret := s.Called(ctx, profileID)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.Profile), ret.Error(1)
}

func (s *ProfileRepositoryMock) GetProfiles(ctx context.Context, f *filter.Filter) (*[]domain.Profile, error) {
	ret := s.Called(ctx, f)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*[]domain.Profile), ret.Error(1)
}

func (s *ProfileRepositoryMock) CreateProfile(ctx context.Context, data *dto.ProfileInputDTO) (*domain.Profile, error) {
	ret := s.Called(ctx, data)
	if ret.Get(0) == nil {
		return nil, ret.Error(1)
	}
	return ret.Get(0).(*domain.Profile), ret.Error(1)
}

func (s *ProfileRepositoryMock) UpdateProfile(ctx context.Context, profile *domain.Profile, data *dto.ProfileInputDTO) error {
	ret := s.Called(ctx, profile, data)
	return ret.Error(0)
}

func (s *ProfileRepositoryMock) DeleteProfiles(ctx context.Context, profileIDs []uint) error {
	ret := s.Called(ctx, profileIDs)
	return ret.Error(0)
}
