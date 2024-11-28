package service

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
)

// NewUserService creates and returns a new instance of domain.UserService with the provided UserRepository.
func NewUserService(r domain.UserRepository) domain.UserService {
	return &userService{
		userRepository: r,
	}
}

// userService provides methods for handling user-related operations such as creating, updating, retrieving, and deleting users.
type userService struct {
	userRepository domain.UserRepository
}

// GenerateUserOutputDTO transforms a domain.User into a dto.UserOutputDTO, mapping essential user and profile details.
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

// GetUserByID retrieves a user by their ID and returns the user's details in a UserOutputDTO structure.
func (s *userService) GetUserByID(ctx context.Context, userID uint) (*dto.UserOutputDTO, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return s.GenerateUserOutputDTO(user), nil
}

// GetUsers retrieves a list of users based on the specified userFilter and returns the users along with pagination info.
// It queries the userRepository to fetch users and the total count, and converts them to UserOutputDTO.
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

// CreateUser creates a new user based on the provided UserInputDTO data and returns the UserOutputDTO or an error.
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

// UpdateUser updates a user with the given userID using the provided UserInputDTO data and returns the updated UserOutputDTO.
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

// DeleteUsers removes users identified by the specified IDs from the repository.
// It returns an error if the operation fails.
func (s *userService) DeleteUsers(ctx context.Context, ids []uint) error {
	return s.userRepository.DeleteUsers(ctx, ids)
}

// ResetUserPassword attempts to reset the password for a user identified by their email.
// If the user does not have a password set, the function returns without changing anything.
// Returns an error if retrieving the user or resetting the password fails.
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

// SetUserPassword sets the password for a user identified by their email, if the user does not already have a password set.
//
// This method retrieves the user by email and checks if a password is already associated with the user. If a password
// exists, it returns an error indicating that the user already has a password. If not, it sets the new password for the user.
//
// Parameters:
//
//	ctx: Context for handling deadlines and cancellations.
//	mail: The email address of the user whose password is to be set.
//	pass: A PasswordInputDTO containing the new password and its confirmation.
//
// Returns:
//
//	An error if the operation fails, or if the user already has a password.
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

// SetUserPhoto updates the photo of a user by their userID with the provided file information in the database.
func (s *userService) SetUserPhoto(ctx context.Context, userID uint, p *domain.File) error {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	return s.userRepository.SetUserPhoto(ctx, user, p)
}

// GenerateUserPhotoURL generates a photo URL for a given user ID by retrieving the user and invoking the repository function.
func (s *userService) GenerateUserPhotoURL(ctx context.Context, userID uint) (string, error) {
	user, err := s.userRepository.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	return s.userRepository.GenerateUserPhotoURL(ctx, user)
}
