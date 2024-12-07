package domain

import (
	"context"
	"github.com/raulaguila/go-api/pkg/utils"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/raulaguila/go-api/pkg/validator"
)

// ProfileTableName represents the name of the database table used for storing user profile information.
const ProfileTableName string = "users_profile"

// Profile represents a user profile with a unique name and associated permissions.
type (
	Profile struct {
		Base
		Name        string        `gorm:"column:name;type:varchar(100);unique;not null;" validate:"required,min=4"`
		Permissions pgutils.JSONB `gorm:"column:permissions;type:jsonb;not null;" validate:"required"`
	}

	ProfileRepository interface {
		CountProfiles(context.Context, *filter.Filter) (int64, error)
		GetProfileByID(context.Context, uint) (*Profile, error)
		GetProfiles(context.Context, *filter.Filter) (*[]Profile, error)
		CreateProfile(context.Context, *Profile) error
		UpdateProfile(context.Context, *Profile) error
		DeleteProfiles(context.Context, []uint) error
	}

	ProfileService interface {
		GenerateProfileOutputDTO(*Profile) *dto.ProfileOutputDTO
		GetProfileByID(context.Context, uint) (*dto.ProfileOutputDTO, error)
		GetProfiles(context.Context, *filter.Filter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error)
		CreateProfile(context.Context, *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error)
		UpdateProfile(context.Context, uint, *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error)
		DeleteProfiles(context.Context, []uint) error
	}
)

// TableName returns the name of the database table associated with the Profile struct.
func (s *Profile) TableName() string {
	return ProfileTableName
}

// ToMap converts the Profile struct into a map with keys as field names and values as corresponding field values.
// It includes the "name" and "permissions" fields from the Profile struct.
func (s *Profile) ToMap() *map[string]any {
	return &map[string]any{
		"name":        s.Name,
		"permissions": s.Permissions,
	}
}

// Bind updates the Profile's Name and Permissions fields based on the provided ProfileInputDTO, and validates the Profile.
func (s *Profile) Bind(p *dto.ProfileInputDTO) error {
	if p != nil {
		s.Name = utils.PointerValue(p.Name, s.Name)
		if p.Permissions != nil {
			s.Permissions = p.Permissions
		}
	}

	return validator.StructValidator.Validate(s)
}
