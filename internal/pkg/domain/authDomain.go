package domain

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/dto"
)

// AuthTableName is the constant string that represents the name of the authentication table in the database schema.
const AuthTableName string = "users_auth"

// Auth represents an authentication entity with credentials and linking to user profile information.
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

// TableName returns the database table name for the Auth struct, which is "users_auth".
func (s *Auth) TableName() string { return AuthTableName }

// ToMap converts the Auth struct into a map with string keys and dynamic value types, representing its fields.
func (s *Auth) ToMap() *map[string]any {
	mapped := map[string]any{
		"status":     s.Status,
		"profile_id": s.ProfileID,
		"token":      nil,
		"password":   nil,
	}

	if s.Token != nil {
		mapped["token"] = *s.Token
		mapped["password"] = *s.Password
	}

	return &mapped
}
