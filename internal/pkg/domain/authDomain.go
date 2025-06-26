package domain

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/dto"
)

const AuthTableName string = "usr_auth"

type (
	Auth struct {
		Base
		Status    bool `gorm:"column:status;"`
		ProfileID uint `gorm:"column:profile_id;" validate:"required,min=1"`
		Profile   *Profile
		Token     *string `gorm:"column:token;"`
		Password  *string `gorm:"column:password;"`
	}

	AuthService interface {
		Login(context.Context, *dto.AuthInputDTO) (*dto.AuthOutputDTO, error)
		Refresh(*User) *dto.AuthOutputDTO
		Me(*User) *dto.UserOutputDTO
	}
)

func (s *Auth) TableName() string { return AuthTableName }

func (s *Auth) ToMap() *map[string]any {
	return &map[string]any{
		"status":     s.Status,
		"profile_id": s.ProfileID,
		"token":      s.Token,
		"password":   s.Password,
	}
}
