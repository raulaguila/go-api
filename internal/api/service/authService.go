package service

import (
	"context"
	"os"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
)

// NewAuthService creates a new instance of the AuthService with a given UserRepository implementation.
func NewAuthService(r domain.UserRepository) domain.AuthService {
	return &authService{
		userRepository: r,
	}
}

// authService is a type that implements authentication-related operations using a UserRepository.
type authService struct {
	userRepository domain.UserRepository
}

// generateUserOutputDTO converts a domain.User to a dto.UserOutputDTO. Returns nil if the input user is nil.
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

// generateAuthOutputDTO generates an AuthOutputDTO containing user data, access token, and refresh token.
// Takes a User object and IP address as input parameters and utilizes environment variables for token expiration.
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

// Login authenticates the user using provided credentials and IP address.
// It returns an AuthOutputDTO on successful authentication or an error if authentication fails.
// Errors could include invalid credentials or if the user account is disabled.
func (s *authService) Login(ctx context.Context, credentials *dto.AuthInputDTO, ip string) (*dto.AuthOutputDTO, error) {
	user, err := s.userRepository.GetUserByMail(ctx, credentials.Login)
	if err != nil {
		return nil, err
	}

	if !user.ValidatePassword(credentials.Password) {
		return nil, myerrors.ErrInvalidCredentials
	}

	if !user.Auth.Status || user.Auth.Password == nil {
		return nil, myerrors.ErrDisabledUser
	}

	return s.generateAuthOutputDTO(user, ip), nil
}

// Me returns a UserOutputDTO for the given domain.User by generating it through the generateUserOutputDTO method.
func (s *authService) Me(user *domain.User) *dto.UserOutputDTO {
	return s.generateUserOutputDTO(user)
}

// Refresh generates a new AuthOutputDTO for the specified user and IP address, renewing access and refresh tokens.
func (s *authService) Refresh(user *domain.User, ip string) *dto.AuthOutputDTO {
	return s.generateAuthOutputDTO(user, ip)
}
