package domain

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/dto"
)

const AuthTableName string = "users_auth"

type (
	Auth struct {
		Base
		Status    bool `gorm:"column:status;type:bool;not null;"`
		ProfileID uint `gorm:"column:profile_id;type:bigint;not null;index;" validate:"required,min=1"`
		Profile   *Profile
		Token     *string `gorm:"column:token;type:varchar(255);unique;index"`
		Password  *string `gorm:"column:password;type:varchar(255);"`
	}

	AuthService interface {
		Login(context.Context, *dto.AuthInputDTO, string) (*dto.AuthOutputDTO, error)
		Refresh(*User, string) *dto.AuthOutputDTO
		Me(*User) *dto.UserOutputDTO
	}
)

func (s *Auth) TableName() string { return AuthTableName }
