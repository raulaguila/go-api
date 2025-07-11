package domain

import (
	"context"
	"crypto/rsa"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/packhub"
	"github.com/raulaguila/go-api/pkg/validator"
)

const UserTableName string = "usr_user"

type (
	User struct {
		BaseInt
		Name     string `gorm:"column:name;" validate:"required,min=5"`
		Username string `gorm:"column:username;" validate:"required,min=5"`
		Email    string `gorm:"column:mail;" validate:"required,email"`
		AuthID   uint   `gorm:"column:auth_id;"`
		Auth     *Auth  `gorm:"constraint:OnDelete:CASCADE"`
	}

	UserRepository interface {
		CountUsers(context.Context, *dto.UserFilter) (int64, error)
		GetUsers(context.Context, *dto.UserFilter) (*[]User, error)
		GetUser(context.Context, *User) error
		GetUserByToken(context.Context, string) (*User, error)
		CreateUser(context.Context, *User) error
		UpdateUser(context.Context, *User) error
		DeleteUsers(context.Context, []uint) error
	}

	UserService interface {
		GenerateUserOutputDTO(*User) *dto.UserOutputDTO
		GetUsers(context.Context, *dto.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error)
		CreateUser(context.Context, *dto.UserInputDTO) (*dto.UserOutputDTO, error)
		UpdateUser(context.Context, uint, *dto.UserInputDTO) (*dto.UserOutputDTO, error)
		DeleteUsers(context.Context, []uint) error
		ResetUserPassword(context.Context, string) error
		SetUserPassword(context.Context, string, *dto.PasswordInputDTO) error
	}
)

func (s *User) AfterDelete(tx *gorm.DB) (err error) {
	return tx.Delete(&Auth{}, s.AuthID).Error
}

func (s *User) TableName() string {
	return UserTableName
}

func (s *User) ToMap() *map[string]any {
	return &map[string]any{
		"name":     s.Name,
		"username": s.Username,
		"mail":     s.Email,
		"auth_id":  s.AuthID,
		"Auth":     *s.Auth.ToMap(),
	}
}

func (s *User) Bind(p *dto.UserInputDTO) error {
	if p != nil {
		s.Name = packhub.PointerValue(p.Name, s.Name)
		s.Username = packhub.PointerValue(p.Username, s.Username)
		s.Email = packhub.PointerValue(p.Email, s.Email)

		s.Auth.Status = packhub.PointerValue(p.Status, s.Auth.Status)
		s.Auth.ProfileID = packhub.PointerValue(p.ProfileID, s.Auth.ProfileID)
	}

	return validator.StructValidator.Validate(s)
}

func (s *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s.Auth.Token = packhub.Pointer(uuid.New().String())
	s.Auth.Password = packhub.Pointer(string(hash))

	return nil
}

func (s *User) ResetPassword() {
	s.Auth.Token = nil
	s.Auth.Password = nil
}

func (s *User) ValidatePassword(password string) bool {
	if s.Auth.Password == nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(*s.Auth.Password), []byte(password)) == nil
}

func (s *User) GenerateToken(expire *time.Duration, parsedToken *rsa.PrivateKey) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"token": s.Auth.Token,
		"iat":   now.Unix(),
	}

	if expire != nil {
		claims["exp"] = now.Add(*expire).Unix()
	}

	return jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(parsedToken)
}
