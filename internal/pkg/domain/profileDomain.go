package domain

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/raulaguila/go-api/pkg/validator"
)

const ProfileTableName string = "users_profile"

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
		CreateProfile(context.Context, *dto.ProfileInputDTO) (*Profile, error)
		UpdateProfile(context.Context, *Profile, *dto.ProfileInputDTO) error
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

func (s *Profile) TableName() string {
	return ProfileTableName
}

func (s *Profile) ToMap() *map[string]any {
	return &map[string]any{
		"name":        s.Name,
		"permissions": s.Permissions,
	}
}

func (s *Profile) Bind(p *dto.ProfileInputDTO) error {
	if p.Name != nil {
		s.Name = *p.Name
	}

	if p.Permissions != nil {
		for key, value := range p.Permissions {
			s.Permissions[key] = value
		}
	}

	return validator.StructValidator.Validate(s)
}
