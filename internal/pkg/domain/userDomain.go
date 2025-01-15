package domain

import (
	"context"
	"crypto/rsa"
	"io"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/utils"
	"github.com/raulaguila/go-api/pkg/validator"
)

const UserTableName string = "users"

type (
	User struct {
		Base
		Name   string `gorm:"column:name;type:varchar(90);not null;" validate:"required,min=5"`
		Email  string `gorm:"column:mail;type:varchar(50);not null;unique;index;" validate:"required,email"`
		AuthID uint   `gorm:"column:auth_id;type:bigint;not null;index;"`
		Auth   *Auth  `gorm:"constraint:OnDelete:CASCADE"`
	}

	File struct {
		Name      string
		Extension string
		File      io.Reader
	}

	UserRepository interface {
		CountUsers(context.Context, *filters.UserFilter) (int64, error)
		GetUsers(context.Context, *filters.UserFilter) (*[]User, error)
		GetUser(context.Context, *User) error
		GetUserByToken(context.Context, string) (*User, error)
		CreateUser(context.Context, *User) error
		UpdateUser(context.Context, *User) error
		DeleteUsers(context.Context, []uint) error
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
	}
)

func (u *User) TableName() string {
	return UserTableName
}

func (u *User) ToMap() *map[string]any {
	mapped := map[string]any{
		"name":    u.Name,
		"mail":    u.Email,
		"auth_id": u.AuthID,
		"Auth": map[string]any{
			"status":     u.Auth.Status,
			"profile_id": u.Auth.ProfileID,
			"token":      nil,
			"password":   nil,
		},
	}

	if u.Auth.Password != nil {
		mapped["Auth"].(map[string]any)["token"] = *u.Auth.Token
		mapped["Auth"].(map[string]any)["password"] = *u.Auth.Password
	}

	return &mapped
}

func (u *User) Bind(p *dto.UserInputDTO) error {
	if p != nil {
		u.Name = utils.PointerValue(p.Name, u.Name)
		u.Email = utils.PointerValue(p.Email, u.Email)
		u.Auth.Status = utils.PointerValue(p.Status, u.Auth.Status)
		u.Auth.ProfileID = utils.PointerValue(p.ProfileID, u.Auth.ProfileID)
	}

	return validator.StructValidator.Validate(u)
}

func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Auth.Token = utils.Pointer(uuid.New().String())
	u.Auth.Password = utils.Pointer(string(hash))

	return nil
}

func (u *User) ResetPassword() {
	u.Auth.Token = nil
	u.Auth.Password = nil
}

func (u *User) ValidatePassword(password string) bool {
	if u.Auth.Password == nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(*u.Auth.Password), []byte(password)) == nil
}

func (u *User) GenerateToken(expire string, parsedToken *rsa.PrivateKey) (string, error) {
	life, err := utils.DurationFromString(expire, time.Minute)
	claims := jwt.MapClaims{
		"token":  u.Auth.Token,
		"expire": err == nil,
		"iat":    time.Now().Unix(),
	}
	if err == nil {
		claims["exp"] = time.Now().Add(life).Unix()
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(parsedToken)
}
