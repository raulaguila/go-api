package domain

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/helper"
	"github.com/raulaguila/go-api/pkg/validator"
)

const UserTableName string = "users"

type (
	User struct {
		Base
		Name      string  `gorm:"column:name;type:varchar(90);not null;" validate:"required,min=5"`
		Email     string  `gorm:"column:mail;type:varchar(50);not null;unique;index;" validate:"required,email"`
		PhotoPath *string `gorm:"column:photo;type:varchar(150);null;default:null;"`
		AuthID    uint    `gorm:"column:auth_id;type:bigint;not null;index;"`
		Auth      *Auth   `gorm:"constraint:OnDelete:CASCADE"`
	}

	File struct {
		Name      string
		Extension string
		File      io.Reader
	}

	UserRepository interface {
		CountUsers(context.Context, *filters.UserFilter) (int64, error)
		GetUsers(context.Context, *filters.UserFilter) (*[]User, error)
		GetUserByID(context.Context, uint) (*User, error)
		GetUserByMail(context.Context, string) (*User, error)
		GetUserByToken(context.Context, string) (*User, error)
		CreateUser(context.Context, *User) error
		UpdateUser(context.Context, *User) error
		DeleteUsers(context.Context, []uint) error
		ResetUserPassword(context.Context, *User) error
		SetUserPassword(context.Context, *User, *dto.PasswordInputDTO) error
		SetUserPhoto(context.Context, *User, *File) error
		GenerateUserPhotoURL(context.Context, *User) (string, error)
	}

	UserService interface {
		GenerateUserOutputDTO(*User) *dto.UserOutputDTO
		GetUserByID(context.Context, uint) (*dto.UserOutputDTO, error)
		GetUsers(context.Context, *filters.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error)
		CreateUser(context.Context, *dto.UserInputDTO) (*dto.UserOutputDTO, error)
		UpdateUser(context.Context, uint, *dto.UserInputDTO) (*dto.UserOutputDTO, error)
		DeleteUsers(context.Context, []uint) error
		ResetUserPassword(context.Context, string) error
		SetUserPassword(context.Context, string, *dto.PasswordInputDTO) error
		SetUserPhoto(context.Context, uint, *File) error
		GenerateUserPhotoURL(context.Context, uint) (string, error)
	}
)

func (u *User) TableName() string {
	return UserTableName
}

func (u *User) ToMap() *map[string]any {
	mapped := &map[string]any{
		"name":    u.Name,
		"mail":    u.Email,
		"auth_id": u.AuthID,
		"photo":   nil,
		"Auth": map[string]any{
			"status":     u.Auth.Status,
			"profile_id": u.Auth.ProfileID,
			"token":      nil,
			"password":   nil,
		},
	}

	if u.PhotoPath != nil {
		(*mapped)["photo"] = *u.PhotoPath
	}

	if u.Auth.Password != nil {
		(*mapped)["Auth"].(map[string]any)["token"] = *u.Auth.Token
		(*mapped)["Auth"].(map[string]any)["password"] = *u.Auth.Password
	}

	return mapped
}

func (u *User) Bind(userDTO *dto.UserInputDTO) error {
	if userDTO != nil {
		if userDTO.Name != nil {
			u.Name = *userDTO.Name
		}
		if userDTO.Email != nil {
			u.Email = *userDTO.Email
		}
		if userDTO.Status != nil {
			u.Auth.Status = *userDTO.Status
		}
		if userDTO.ProfileID != nil {
			u.Auth.ProfileID = *userDTO.ProfileID
		}
	}

	return validator.StructValidator.Validate(u)
}

func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(*u.Auth.Password), []byte(password)) == nil
}

func (u *User) GenerateToken(expire, originalKey, ip string) (string, error) {
	decodedKey, err := base64.StdEncoding.DecodeString(originalKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %v", err.Error())
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedKey)
	if err != nil {
		return "", fmt.Errorf("could not parse key: %v", err.Error())
	}

	now := time.Now()
	claims := jwt.MapClaims{"token": u.Auth.Token, "ip": ip, "iat": now.Unix()}

	life, err := helper.DurationFromString(expire, time.Minute)
	if err == nil {
		claims["exp"] = now.Add(life).Unix()
	}
	claims["expire"] = err == nil

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}
