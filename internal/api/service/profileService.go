package service

import (
	"context"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/raulaguila/go-api/pkg/utils"

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
		ID:          &profile.ID,
		Name:        &profile.Name,
		Permissions: &profile.Permissions,
	}
}

func (s *profileService) GetProfileByID(ctx context.Context, profileID uint) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.GetProfileByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) GetProfiles(ctx context.Context, profileFilter *filter.Filter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error) {
	profiles, err := s.profileRepository.GetProfiles(ctx, profileFilter)
	if err != nil {
		return nil, err
	}

	count, err := s.profileRepository.CountProfiles(ctx, profileFilter)
	if err != nil {
		return nil, err
	}

	outputProfiles := make([]dto.ProfileOutputDTO, 0)
	for _, profile := range *profiles {
		outputProfiles = append(outputProfiles, *s.GenerateProfileOutputDTO(&profile))
	}

	return &dto.ItemsOutputDTO[dto.ProfileOutputDTO]{
		Items: outputProfiles,
		Pagination: dto.PaginationDTO{
			CurrentPage: uint(utils.Max(profileFilter.Page, 1)),
			PageSize:    uint(utils.Max(profileFilter.Limit, len(outputProfiles))),
			TotalItems:  uint(count),
			TotalPages:  uint(profileFilter.CalcPages(count)),
		},
	}, nil
}

func (s *profileService) CreateProfile(ctx context.Context, data *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{Permissions: pgutils.JSONB{}}
	if err := profile.Bind(data); err != nil {
		return nil, err
	}

	if err := s.profileRepository.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) UpdateProfile(ctx context.Context, profileID uint, data *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.GetProfileByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	if err := profile.Bind(data); err != nil {
		return nil, err
	}

	if err := s.profileRepository.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) DeleteProfiles(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return s.profileRepository.DeleteProfiles(ctx, ids)
}
