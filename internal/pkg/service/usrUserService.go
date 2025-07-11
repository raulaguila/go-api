package service

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/packhub"
	"github.com/raulaguila/go-api/pkg/utils"
)

func NewUserService(r domain.UserRepository) domain.UserService {
	return &userService{
		repository: r,
	}
}

type userService struct {
	repository domain.UserRepository
}

func (s *userService) GenerateUserOutputDTO(user *domain.User) *dto.UserOutputDTO {
	return &dto.UserOutputDTO{
		ID:       &user.ID,
		Name:     &user.Name,
		Username: &user.Username,
		Email:    &user.Email,
		Status:   &user.Auth.Status,
		New:      packhub.Pointer(user.Auth.Password == nil),
		Profile: &dto.ProfileOutputDTO{
			ID:   &user.Auth.Profile.ID,
			Name: &user.Auth.Profile.Name,
		},
	}
}

func (s *userService) GetUserByID(ctx context.Context, userID uint) (*dto.UserOutputDTO, error) {
	user := &domain.User{BaseInt: domain.BaseInt{ID: userID}}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

func (s *userService) GetUsers(ctx context.Context, userFilter *dto.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error) {
	users, err := s.repository.GetUsers(ctx, userFilter)
	if err != nil {
		return nil, err
	}

	count, err := s.repository.CountUsers(ctx, userFilter)
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

	if err := s.repository.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	user = &domain.User{BaseInt: domain.BaseInt{ID: user.ID}}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

func (s *userService) UpdateUser(ctx context.Context, userID uint, data *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user := &domain.User{BaseInt: domain.BaseInt{ID: userID}}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	if err := user.Bind(data); err != nil {
		return nil, err
	}

	if err := s.repository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	user = &domain.User{BaseInt: domain.BaseInt{ID: userID}}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

func (s *userService) DeleteUsers(ctx context.Context, ids []uint) error {
	return s.repository.DeleteUsers(ctx, ids)
}

func (s *userService) ResetUserPassword(ctx context.Context, mail string) error {
	user := &domain.User{Email: mail}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return err
	}

	if user.Auth.Password == nil && user.Auth.Token == nil {
		return nil
	}

	user.ResetPassword()
	return s.repository.UpdateUser(ctx, user)
}

func (s *userService) SetUserPassword(ctx context.Context, mail string, pass *dto.PasswordInputDTO) error {
	user := &domain.User{Email: mail}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return err
	}

	if user.Auth.Password != nil {
		return utils.ErrUserHasPass
	}

	if err := user.SetPassword(*pass.Password); err != nil {
		return err
	}

	return s.repository.UpdateUser(ctx, user)
}
