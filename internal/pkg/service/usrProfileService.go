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

func (s *profileService) GenerateProfileOutputDTO(output *domain.Profile) *dto.ProfileOutputDTO {
	return &dto.ProfileOutputDTO{
		ID:   &output.ID,
		Name: &output.Name,
		Permissions: func() *pq.StringArray {
			if output.Permissions != nil {
				return &output.Permissions
			}
			return nil
		}(),
	}
}

func (s *profileService) GetProfileByID(ctx context.Context, id uint) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{BaseInt: domain.BaseInt{ID: id}}
	if err := s.repository.GetProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) GetProfiles(ctx context.Context, f *dto.ProfileFilter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error) {
	profiles, err := s.repository.GetProfiles(ctx, f)
	if err != nil {
		return nil, err
	}

	count, err := s.repository.CountProfiles(ctx, f)
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
			Page:       uint(packhub.Max(f.Page, 1)),
			Limit:      uint(packhub.Max(f.Limit, len(outputProfiles))),
			TotalItems: uint(count),
			TotalPages: uint(f.CalcPages(count)),
		},
	}, nil
}

func (s *profileService) ListProfiles(ctx context.Context, f *dto.ProfileFilter) (*[]dto.ItemOutputDTO, error) {
	f.Page = 0
	f.Limit = 0

	profiles, err := s.repository.GetProfiles(ctx, f)
	if err != nil {
		return nil, err
	}

	outputProfiles := make([]dto.ItemOutputDTO, len(*profiles))
	for i, profile := range *profiles {
		outputProfiles[i] = dto.ItemOutputDTO{
			ID:   &profile.ID,
			Name: &profile.Name,
		}
	}

	return &outputProfiles, nil
}

func (s *profileService) CreateProfile(ctx context.Context, input *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{Permissions: []string{}}
	if err := profile.Bind(input); err != nil {
		return nil, err
	}

	if err := s.repository.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

func (s *profileService) UpdateProfile(ctx context.Context, id uint, input *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{BaseInt: domain.BaseInt{ID: id}}
	if err := s.repository.GetProfile(ctx, profile); err != nil {
		return nil, err
	}

	if err := profile.Bind(input); err != nil {
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
