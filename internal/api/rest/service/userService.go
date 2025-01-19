package service

import (
	"context"
	"fmt"

	"github.com/raulaguila/packhub"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/utils"
)

func NewUserService(r domain.UserRepository) domain.UserService {
	return &userService{
		userRepository: r,
	}
}

type userService struct {
	userRepository domain.UserRepository
}

func (s *userService) GenerateUserOutputDTO(user *domain.User) *dto.UserOutputDTO {
	return &dto.UserOutputDTO{
		ID:     &user.ID,
		Name:   &user.Name,
		Email:  &user.Email,
		Status: &user.Auth.Status,
		Profile: &dto.ProfileOutputDTO{
			ID:          &user.Auth.Profile.ID,
			Name:        &user.Auth.Profile.Name,
			Permissions: &user.Auth.Profile.Permissions,
		},
	}
}

func (s *userService) GetUserByID(ctx context.Context, userID uint) (*dto.UserOutputDTO, error) {
	user := &domain.User{Base: domain.Base{ID: userID}}
	if err := s.userRepository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

func (s *userService) GetUsers(ctx context.Context, userFilter *filters.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error) {
	users, err := s.userRepository.GetUsers(ctx, userFilter)
	if err != nil {
		fmt.Printf("GetUsers Error: %v\n", err)
		return nil, err
	}

	count, err := s.userRepository.CountUsers(ctx, userFilter)
	if err != nil {
		fmt.Printf("CountUsers Error: %v\n", err)
		return nil, err
	}

	outputUsers := make([]dto.UserOutputDTO, 0)
	for _, user := range *users {
		outputUsers = append(outputUsers, *s.GenerateUserOutputDTO(&user))
	}

	return &dto.ItemsOutputDTO[dto.UserOutputDTO]{
		Items: outputUsers,
		Pagination: dto.PaginationDTO{
			CurrentPage: uint(packhub.Max(userFilter.Page, 1)),
			PageSize:    uint(packhub.Max(userFilter.Limit, len(outputUsers))),
			TotalItems:  uint(count),
			TotalPages:  uint(userFilter.CalcPages(count)),
		},
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, data *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user := &domain.User{Auth: &domain.Auth{}}
	if err := user.Bind(data); err != nil {
		return nil, err
	}

	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	user = &domain.User{Base: domain.Base{ID: user.ID}}
	if err := s.userRepository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uint, data *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user := &domain.User{Base: domain.Base{ID: userID}}
	if err := s.userRepository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	if err := user.Bind(data); err != nil {
		return nil, err
	}

	if err := s.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	user = &domain.User{Base: domain.Base{ID: userID}}
	if err := s.userRepository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

func (s *userService) DeleteUsers(ctx context.Context, ids []uint) error {
	return s.userRepository.DeleteUsers(ctx, ids)
}

func (s *userService) ResetUserPassword(ctx context.Context, mail string) error {
	user := &domain.User{Email: mail}
	if err := s.userRepository.GetUser(ctx, user); err != nil {
		return err
	}

	if user.Auth.Password == nil && user.Auth.Token == nil {
		return nil
	}

	user.ResetPassword()
	return s.userRepository.UpdateUser(ctx, user)
}

func (s *userService) SetUserPassword(ctx context.Context, mail string, pass *dto.PasswordInputDTO) error {
	user := &domain.User{Email: mail}
	if err := s.userRepository.GetUser(ctx, user); err != nil {
		return err
	}

	if user.Auth.Password != nil {
		return utils.ErrUserHasPass
	}

	if err := user.SetPassword(*pass.Password); err != nil {
		return err
	}

	return s.userRepository.UpdateUser(ctx, user)
}
