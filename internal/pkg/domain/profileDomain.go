package domain

import (
	"context"

	"github.com/lib/pq"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/utils"
	"github.com/raulaguila/go-api/pkg/validator"
)

const ProfileTableName string = "users_profile"

type (
	Profile struct {
		Base
		Name        string         `gorm:"column:name;type:varchar(100);unique;not null;" validate:"required,min=4"`
		Permissions pq.StringArray `gorm:"column:permissions;type:text[];not null;" validate:"required"`
	}

	ProfileRepository interface {
		CountProfiles(ctx context.Context, f *filter.Filter) (int64, error)
		GetProfile(ctx context.Context, p *Profile) error
		GetProfiles(ctx context.Context, f *filter.Filter) (*[]Profile, error)
		CreateProfile(ctx context.Context, p *Profile) error
		UpdateProfile(ctx context.Context, p *Profile) error
		DeleteProfiles(ctx context.Context, i []uint) error
	}

	ProfileService interface {
		GenerateProfileOutputDTO(p *Profile) *dto.ProfileOutputDTO
		GetProfileByID(ctx context.Context, id uint) (*dto.ProfileOutputDTO, error)
		GetProfiles(ctx context.Context, f *filter.Filter) (*dto.ItemsOutputDTO[dto.ProfileOutputDTO], error)
		CreateProfile(ctx context.Context, pdto *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error)
		UpdateProfile(ctx context.Context, id uint, pdto *dto.ProfileInputDTO) (*dto.ProfileOutputDTO, error)
		DeleteProfiles(ctx context.Context, ids []uint) error
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
	if p != nil {
		s.Name = utils.PointerValue(p.Name, s.Name)
		s.Permissions = utils.PointerValue(p.Permissions, s.Permissions)
	}

	return validator.StructValidator.Validate(s)
}
