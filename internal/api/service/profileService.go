package service

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

func NewProfileService(r domain.ProfileRepository) domain.ProfileService {
	return &profileService{
		profileRepository: r,
	}
}

type profileService struct {
	profileRepository domain.ProfileRepository
}

func (s *profileService) GenerateProfileOutputDTO(profile *domain.Profile) *dto.ProfileOutputDTO {
	return &dto.ProfileOutputDTO{
		Id:          &profile.Id,
		Name:        &profile.Name,
		Permissions: profile.Permissions,
	}
}

// GetProfileByID Implementation of 'GetProfileByID'.
func (s *profileService) GetProfileByID(ctx context.Context, profileID uint) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.GetProfileByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

// GetProfiles Implementation of 'GetProfiles'.
func (s *profileService) GetProfiles(ctx context.Context, filter *filter.Filter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error) {
	profiles, err := s.profileRepository.GetProfiles(ctx, filter)
	if err != nil {
		return nil, err
	}

	count, err := s.profileRepository.CountProfiles(ctx, filter)
	if err != nil {
		return nil, err
	}

	outputProfiles := make([]dto.ProfileOutputDTO, 0)
	for _, profile := range *profiles {
		outputProfiles = append(outputProfiles, *s.GenerateProfileOutputDTO(&profile))
	}
	return &dto.ItemsOutputDTO[dto.ProfileOutputDTO]{
		Count: count,
		Items: outputProfiles,
		Pages: filter.CalcPages(count),
	}, nil
}

// CreateProfile Implementation of 'CreateProfile'.
func (s *profileService) CreateProfile(ctx context.Context, data *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.CreateProfile(ctx, data)
	if err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

// UpdateProfile Implementation of 'UpdateProfile'.
func (s *profileService) UpdateProfile(ctx context.Context, profileID uint, data *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.GetProfileByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	if err := s.profileRepository.UpdateProfile(ctx, profile, data); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

// DeleteProfiles Implementation of 'DeleteProfiles'.
func (s *profileService) DeleteProfiles(ctx context.Context, ids []uint) error {
	return s.profileRepository.DeleteProfiles(ctx, ids)
}
