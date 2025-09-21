package service

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/erro"
	"github.com/raulaguila/go-api/pkg/packhub"
)

func NewUserService(r domain.UserRepository) domain.UserService {
	return &userService{
		repository: r,
	}
}

type userService struct {
	repository domain.UserRepository
}

func (s *userService) GenerateUserOutputDTO(output *domain.User) *dto.UserOutputDTO {
	return &dto.UserOutputDTO{
		ID:       &output.ID,
		Name:     &output.Name,
		Username: &output.Username,
		Email:    &output.Email,
		Status:   &output.Auth.Status,
		New:      packhub.Pointer(output.Auth.Password == nil),
		Profile: &dto.ProfileOutputDTO{
			ID:   &output.Auth.Profile.ID,
			Name: &output.Auth.Profile.Name,
		},
	}
}

func (s *userService) GetUsers(ctx context.Context, f *dto.UserFilter) (*dto.ItemsOutputDTO[dto.UserOutputDTO], error) {
	users, err := s.repository.GetUsers(ctx, f)
	if err != nil {
		return nil, err
	}

	count, err := s.repository.CountUsers(ctx, f)
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
			Page:       uint(packhub.Max(f.Page, 1)),
			Limit:      uint(packhub.Max(f.Limit, len(outputUsers))),
			TotalItems: uint(count),
			TotalPages: uint(f.CalcPages(count)),
		},
	}, nil
}

func (s *userService) CreateUser(ctx context.Context, input *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user := &domain.User{Auth: &domain.Auth{}}
	if err := user.Bind(input); err != nil {
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

func (s *userService) UpdateUser(ctx context.Context, id uint, input *dto.UserInputDTO) (*dto.UserOutputDTO, error) {
	user := &domain.User{BaseInt: domain.BaseInt{ID: id}}
	if err := s.repository.GetUser(ctx, user); err != nil {
		return nil, err
	}

	if err := user.Bind(input); err != nil {
		return nil, err
	}

	if err := s.repository.UpdateUser(ctx, user); err != nil {
		return nil, err
	}

	user = &domain.User{BaseInt: domain.BaseInt{ID: id}}
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
		return erro.ErrUserHasPass
	}

	if err := user.SetPassword(*pass.Password); err != nil {
		return err
	}

	return s.repository.UpdateUser(ctx, user)
}
