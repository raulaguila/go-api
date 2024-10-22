package service

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
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
			ID:   &user.Auth.Profile.ID,
			Name: &user.Auth.Profile.Name,
		},
	}
}

// GetUserByID Implementation of 'GetUserByID'.
func (s *userService) GetUserByID(ctx context.Context, userID uint) (*dto.UserOutputDTO, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

// GetUsers Implementation of 'GetUsers'.
func (s *userService) GetUsers(ctx context.Context, userFilter *filters.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error) {
	users, err := s.userRepository.GetUsers(ctx, userFilter)
	if err != nil {
		return nil, err
	}

	count, err := s.userRepository.CountUsers(ctx, userFilter)
	if err != nil {
		return nil, err
	}

	outputUsers := make([]dto.UserOutputDTO, 0)
	for _, user := range *users {
		outputUsers = append(outputUsers, *s.GenerateUserOutputDTO(&user))
	}

	return &dto.ItemsOutputDTO[dto.UserOutputDTO]{
		Items: outputUsers,
		Pagination: dto.PaginationDTO{
			CurrentPage: uint(max(userFilter.Page, 1)),
			PageSize:    uint(max(userFilter.Limit, len(outputUsers))),
			TotalItems:  uint(count),
			TotalPages:  uint(userFilter.CalcPages(count)),
		},
	}, nil
}

// CreateUser Implementation of 'CreateUser'.
func (s *userService) CreateUser(ctx context.Context, data *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user := &domain.User{Auth: &domain.Auth{}}
	if err := user.Bind(data); err != nil {
		return nil, err
	}

	if err := s.userRepository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	user, err := s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

// UpdateUser Implementation of 'UpdateUser'.
func (s *userService) UpdateUser(ctx context.Context, userID uint, data *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if err := user.Bind(data); err != nil {
		return nil, err
	}

	if err := s.userRepository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	user, err = s.userRepository.GetUserByID(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

// DeleteUsers Implementation of 'DeleteUsers'.
func (s *userService) DeleteUsers(ctx context.Context, ids []uint) error {
	return s.userRepository.DeleteUsers(ctx, ids)
}

// ResetUserPassword Implementation of 'ResetUserPassword'.
func (s *userService) ResetUserPassword(ctx context.Context, mail string) error {
	user, err := s.userRepository.GetUserByMail(ctx, mail)
	if err != nil {
		return err
	}

	if user.Auth.Password == nil {
		return nil
	}

	return s.userRepository.ResetUserPassword(ctx, user)
}

// SetUserPassword Implementation of 'SetUserPassword'.
func (s *userService) SetUserPassword(ctx context.Context, mail string, pass *dto.PasswordInputDTO) error {
	user, err := s.userRepository.GetUserByMail(ctx, mail)
	if err != nil {
		return err
	}

	if user.Auth.Password != nil {
		return myerrors.ErrUserHasPass
	}

	return s.userRepository.SetUserPassword(ctx, user, pass)
}

// SetUserPhoto Implementation of 'SetUserPhoto'.
func (s *userService) SetUserPhoto(ctx context.Context, userID uint, p *domain.File) error {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.userRepository.SetUserPhoto(ctx, user, p)
}

func (s *userService) GenerateUserPhotoURL(ctx context.Context, userID uint) (string, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	return s.userRepository.GenerateUserPhotoURL(ctx, user)
}
