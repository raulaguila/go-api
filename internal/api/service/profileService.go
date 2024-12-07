package service

import (
	"context"
	"github.com/raulaguila/go-api/pkg/utils"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
)

// NewProfileService initializes a new instance of ProfileService with the provided ProfileRepository.
func NewProfileService(r domain.ProfileRepository) domain.ProfileService {
	return &profileService{
		profileRepository: r,
	}
}

// profileService handles business logic related to user profiles by utilizing a ProfileRepository for data operations.
type profileService struct {
	profileRepository domain.ProfileRepository
}

// GenerateProfileOutputDTO transforms a domain.Profile into a dto.ProfileOutputDTO for API output purposes.
func (s *profileService) GenerateProfileOutputDTO(profile *domain.Profile) *dto.ProfileOutputDTO {
	return &dto.ProfileOutputDTO{
		ID:          &profile.ID,
		Name:        &profile.Name,
		Permissions: profile.Permissions,
	}
}

// GetProfileByID retrieves a user profile by its unique identifier and returns a ProfileOutputDTO or an error.
func (s *profileService) GetProfileByID(ctx context.Context, profileID uint) (*dto.ProfileOutputDTO, error) {
	profile, err := s.profileRepository.GetProfileByID(ctx, profileID)
	if err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

// GetProfiles retrieves a list of profiles based on the provided filter criteria.
// It returns the profiles along with pagination details or an error if the operation fails.
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

// CreateProfile creates a new user profile in the system using the provided ProfileInputDTO data.
func (s *profileService) CreateProfile(ctx context.Context, data *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error) {
	profile := &domain.Profile{Permissions: map[string]any{}}
	if err := profile.Bind(data); err != nil {
		return nil, err
	}

	if err := s.profileRepository.CreateProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.GenerateProfileOutputDTO(profile), nil
}

// UpdateProfile updates an existing profile with the specified profileID using the provided ProfileInputDTO data.
// Returns a ProfileOutputDTO containing the updated profile information or an error if the update operation fails.
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

// DeleteProfiles removes the user profiles with the specified IDs from the repository.
// It returns an error if the deletion operation fails. If no IDs are provided, it does nothing and returns nil.
func (s *profileService) DeleteProfiles(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}

	return s.profileRepository.DeleteProfiles(ctx, ids)
}
