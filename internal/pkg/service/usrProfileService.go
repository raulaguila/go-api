package service

import (
	"context"

	"github.com/lib/pq"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/packhub"
)

func NewProfileService(r domain.ProfileRepository) domain.ProfileService {
	return &profileService{
		repository: r,
	}
}

type profileService struct {
	repository domain.ProfileRepository
}

func (s *profileService) GenerateProfileOutputDTO(profile *domain.Profile) *dto.ProfileOutputDTO {
	return &dto.ProfileOutputDTO{
		ID:   &profile.ID,
		Name: &profile.Name,
		Permissions: func() *pq.StringArray {
			if profile.Permissions != nil {
				return &profile.Permissions
			}
			return nil
		}(),
	}

	//if profile.Permissions != nil {
	//	result.Permissions = &profile.Permissions
	//}
	//
	//return result
}

func (s *profileService) GetProfileByID(ctx context.Context, profileID uint) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{BaseInt: domain.BaseInt{ID: profileID}}
	if err := s.repository.GetProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) GetProfiles(ctx context.Context, profileFilter *dto.ProfileFilter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error) {
	profiles, err := s.repository.GetProfiles(ctx, profileFilter)
	if err != nil {
		return nil, err
	}

	count, err := s.repository.CountProfiles(ctx, profileFilter)
	if err != nil {
		return nil, err
	}

	outputProfiles := make([]dto.ProfileOutputDTO, len(*profiles))
	for i, profile := range *profiles {
		outputProfiles[i] = *s.GenerateProfileOutputDTO(&profile)
	}

	return &dto.ItemsOutputDTO[dto.ProfileOutputDTO]{
		Items: outputProfiles,
		Pagination: dto.PaginationDTO{
			CurrentPage: uint(packhub.Max(profileFilter.Page, 1)),
			PageSize:    uint(packhub.Max(profileFilter.Limit, len(outputProfiles))),
			TotalItems:  uint(count),
			TotalPages:  uint(profileFilter.CalcPages(count)),
		},
	}, nil
}

func (s *profileService) CreateProfile(ctx context.Context, pdto *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{Permissions: []string{}}
	if err := profile.Bind(pdto); err != nil {
		return nil, err
	}

	if err := s.repository.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) UpdateProfile(ctx context.Context, id uint, pdto *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{BaseInt: domain.BaseInt{ID: id}}
	if err := s.repository.GetProfile(ctx, profile); err != nil {
		return nil, err
	}

	if err := profile.Bind(pdto); err != nil {
		return nil, err
	}

	if err := s.repository.UpdateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) DeleteProfiles(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return s.repository.DeleteProfiles(ctx, ids)
}
