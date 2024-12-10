package service

import (
	"context"
	"os"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/utils"
)

func NewAuthService(r domain.UserRepository) domain.AuthService {
	return &authService{
		userRepository: r,
	}
}

type authService struct {
	userRepository domain.UserRepository
}

func (s *authService) generateUserOutputDTO(user *domain.User) *dto.UserOutputDTO {
	if user == nil {
		return nil
	}

	return &dto.UserOutputDTO{
		ID:     &user.ID,
		Name:   &user.Name,
		Email:  &user.Email,
		Status: &user.Auth.Status,
		Profile: &dto.ProfileOutputDTO{
			ID:          &user.Auth.Profile.ID,
			Name:        &user.Auth.Profile.Name,
			Permissions: user.Auth.Profile.Permissions,
		},
	}
}

func (s *authService) generateAuthOutputDTO(user *domain.User, ip string) *dto.AuthOutputDTO {
	accessTime := os.Getenv("ACCESS_TOKEN_EXPIRE")
	refreshTime := os.Getenv("RFRESH_TOKEN_EXPIRE")

	accessToken, _ := user.GenerateToken(accessTime, os.Getenv("ACCESS_TOKEN_PRIVAT"), ip)
	refreshToken, _ := user.GenerateToken(refreshTime, os.Getenv("RFRESH_TOKEN_PRIVAT"), ip)

	return &dto.AuthOutputDTO{
		User:         s.generateUserOutputDTO(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (s *authService) Login(ctx context.Context, credentials *dto.AuthInputDTO, ip string) (*dto.AuthOutputDTO, error) {
	user, err := s.userRepository.GetUserByMail(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(credentials.Password) {
		return nil, utils.ErrInvalidCredentials
	}

	if !user.Auth.Status || user.Auth.Password == nil {
		return nil, utils.ErrDisabledUser
	}

	return s.generateAuthOutputDTO(user, ip), nil
}

func (s *authService) Me(user *domain.User) *dto.UserOutputDTO {
	return s.generateUserOutputDTO(user)
}

func (s *authService) Refresh(user *domain.User, ip string) *dto.AuthOutputDTO {
	return s.generateAuthOutputDTO(user, ip)
}
