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
	"github.com/raulaguila/go-api/pkg/utils"
	"github.com/raulaguila/go-api/pkg/validator"
)

// UserTableName specifies the database table name for storing User entities.
const UserTableName string = "users"

// User represents a user entity containing basic information and associated authentication details.
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

// TableName returns the name of the database table associated with the User struct.
func (u *User) TableName() string {
	return UserTableName
}

// ToMap converts the User struct into a map with string keys and dynamic value types, representing its fields.
func (u *User) ToMap() *map[string]any {
	mapped := map[string]any{
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
		mapped["photo"] = *u.PhotoPath
	}

	if u.Auth.Password != nil {
		mapped["Auth"].(map[string]any)["token"] = *u.Auth.Token
		mapped["Auth"].(map[string]any)["password"] = *u.Auth.Password
	}

	return &mapped
}

// Bind updates the User fields with values from the provided UserInputDTO if they are not nil.
// Returns an error if the updated User does not pass validation.
func (u *User) Bind(p *dto.UserInputDTO) error {
	if p != nil {
		u.Name = utils.PointerValue(p.Name, u.Name)
		u.Email = utils.PointerValue(p.Email, u.Email)
		u.Auth.Status = utils.PointerValue(p.Status, u.Auth.Status)
		u.Auth.ProfileID = utils.PointerValue(p.ProfileID, u.Auth.ProfileID)
	}

	return validator.StructValidator.Validate(u)
}

// ValidatePassword compares the provided password with the stored hashed password for the user.
// Returns true if the password matches, otherwise false.
func (u *User) ValidatePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(*u.Auth.Password), []byte(password)) == nil
}

// GenerateToken creates a JWT token for a user with a specified expiration, private key, and IP address.
// It returns the signed token string or an error if the process fails.
// The expiration duration is parsed from a string format and used to set the token's expiration claim.
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

	life, err := utils.DurationFromString(expire, time.Minute)
	if err == nil {
		claims["exp"] = now.Add(life).Unix()
	}
	claims["expire"] = err == nil

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(key)
}
