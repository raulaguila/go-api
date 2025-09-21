package service

import (
	"context"
	"time"

	"github.com/raulaguila/go-api/configs"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/erro"
)

func NewAuthService(r domain.UserRepository) domain.AuthService {
	return &authService{
		repository: r,
	}
}

type authService struct {
	repository domain.UserRepository
}

func (s *authService) generateUserOutputDTO(output *domain.User) *dto.UserOutputDTO {
	if output == nil {
		return nil
	}

	return &dto.UserOutputDTO{
		ID:       &output.ID,
		Name:     &output.Name,
		Username: &output.Username,
		Email:    &output.Email,
		Status:   &output.Auth.Status,
		Profile: &dto.ProfileOutputDTO{
			ID:          &output.Auth.Profile.ID,
			Name:        &output.Auth.Profile.Name,
			Permissions: &output.Auth.Profile.Permissions,
		},
	}
}

func (s *authService) generateAuthOutputDTO(user *domain.User, expiration bool) *dto.AuthOutputDTO {
	accessToken, _ := user.GenerateToken(func() *time.Duration {
		if expiration {
			return &configs.AccessExpiration
		}
		return nil
	}(), configs.AccessPrivateKey)
	refreshToken, _ := user.GenerateToken(func() *time.Duration {
		if expiration {
			return &configs.RefreshExpiration
		}
		return nil
	}(), configs.RefreshPrivateKey)

	return &dto.AuthOutputDTO{
		User:         s.generateUserOutputDTO(user),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func (s *authService) Login(ctx context.Context, input *dto.AuthInputDTO) (*dto.AuthOutputDTO, error) {
	user := &domain.User{Username: input.Login}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	if !user.ValidatePassword(input.Password) {
		return nil, erro.ErrInvalidCredentials
	}

	if !user.Auth.Status || user.Auth.Password == nil {
		return nil, erro.ErrDisabledUser
	}

	return s.generateAuthOutputDTO(user, input.Expiration), nil
}

func (s *authService) Me(user *domain.User) *dto.UserOutputDTO {
	return s.generateUserOutputDTO(user)
}

func (s *authService) Refresh(user *domain.User, expiration bool) *dto.AuthOutputDTO {
	return s.generateAuthOutputDTO(user, expiration)
}
