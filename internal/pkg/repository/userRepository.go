package repository

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/internal/pkg/myerrors"
	"github.com/raulaguila/go-api/internal/pkg/postgre"
	"github.com/raulaguila/go-api/pkg/minioutils"
)

// NewUserRepository creates and returns a new instance of UserRepository using the provided gorm DB and Minio client.
func NewUserRepository(db *gorm.DB, minioClient *minioutils.Minio) domain.UserRepository {
	return &userRepository{
		db:    db,
		minio: minioClient,
	}
}

// userRepository is a struct that provides methods to interact with the users in the database.
type userRepository struct {
	db    *gorm.DB
	minio *minioutils.Minio
}

// applyFilter applies the provided filters to the user database query within a context and returns a modified *gorm.DB.
// It handles joins between user, auth, and profile tables, applies search filters, and sets the order and grouping.
func (s *userRepository) applyFilter(ctx context.Context, filter *filters.UserFilter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if filter != nil {
		if filter.ProfileID != 0 {
			db = db.Where(domain.AuthTableName+".profile_id = ?", filter.ProfileID)
		}
		db = db.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.auth_id", domain.AuthTableName, domain.AuthTableName, domain.UserTableName))
		db = db.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.profile_id", domain.ProfileTableName, domain.ProfileTableName, domain.AuthTableName))
		db = filter.ApplySearchLike(db,
			domain.UserTableName+".name",
			domain.UserTableName+".mail",
			domain.ProfileTableName+".name",
		)
		tbName := domain.UserTableName
		db = filter.ApplyOrder(db, &tbName)
	}

	return db.Group(domain.UserTableName + ".id")
}

// CountUsers returns the total number of users based on the provided filter criteria.
// It uses applyFilter to refine the selection of users before counting.
// Takes a context for managing request-scoped values and cancellation signals.
// Returns an integer count of users and an error if one occurred during the operation.
func (s *userRepository) CountUsers(ctx context.Context, filter *filters.UserFilter) (int64, error) {
	var count int64

	db := s.applyFilter(ctx, filter)
	return count, db.Model(&domain.User{}).Count(&count).Error
}

// GetUsers retrieves a list of users from the database based on the provided filter criteria.
// The filter may include parameters for pagination and search. If the filter is not nil, pagination is applied.
// Returns a pointer to a slice of User objects and any error encountered during the database query.
func (s *userRepository) GetUsers(ctx context.Context, filter *filters.UserFilter) (*[]domain.User, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	users := &[]domain.User{}
	return users, db.Preload(postgre.AuthProfile).Find(users).Error
}

// GetUserByID retrieves a user from the database by their unique user ID.
func (s *userRepository) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	user := &domain.User{}
	return user, s.db.WithContext(ctx).Preload(postgre.AuthProfile).First(user, userID).Error
}

// GetUserByMail retrieves a user from the database using the provided email address.
func (s *userRepository) GetUserByMail(ctx context.Context, mail string) (*domain.User, error) {
	user := new(domain.User)
	return user, s.db.WithContext(ctx).Preload(postgre.AuthProfile).First(user, "mail = ?", mail).Error
}

// GetUserByToken retrieves a user by their authentication token. It returns a User object or an error if not found.
func (s *userRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	auth := new(domain.Auth)
	if err := s.db.WithContext(ctx).First(auth, "token = ?", token).Error; err != nil {
		return nil, err
	}

	user := new(domain.User)
	return user, s.db.WithContext(ctx).Preload(postgre.AuthProfile).First(user, "auth_id = ?", auth.ID).Error
}

// CreateUser adds a new user record into the database with its associated data.
// It utilizes a gorm session to ensure all associations are fully saved.
// The function requires a context for managing request-scoped values and cancellation signals,
// and a pointer to a domain.User struct that represents the user details to be saved.
// Returns an error if the creation fails.
func (s *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).Create(user).Error
}

// UpdateUser updates the user and their associated authentication details in the database within a transaction.
// It takes a context and a user object containing the new user information.
// Returns an error if the operation fails at any point.
func (s *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.db.Model(user.Auth).Updates(user.Auth.ToMap()).Error; err != nil {
			return err
		}

		return s.db.Model(user).Updates(user.ToMap()).Error
	})
}

// DeleteUsers removes users and their associated authentication records from the database using a transaction.
// It requires a context for database operations and a slice of user IDs to delete.
// Returns an error if any database operation fails.
func (s *userRepository) DeleteUsers(ctx context.Context, toDelete []uint) error {
	if len(toDelete) == 0 {
		return nil
	}

	users := make([]*domain.User, len(toDelete))
	if err := s.db.WithContext(ctx).Find(&users, toDelete).Error; err != nil {
		return err
	}

	auths := make([]uint, len(toDelete))
	for _, user := range users {
		auths = append(auths, user.AuthID)
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Delete(users).Error; err != nil {
			return err
		}

		return tx.WithContext(ctx).Delete(&domain.Auth{}, auths).Error
	})
}

// ResetUserPassword resets a user's password by clearing the Password and Token fields in the associated Auth object.
// It then updates the user with these changes in the database.
// Params:
// - ctx: The context for request-scoped values.
// - user: A pointer to the User object whose password needs to be reset.
// Returns an error if the update operation fails.
func (s *userRepository) ResetUserPassword(ctx context.Context, user *domain.User) error {
	user.Auth.Password = nil
	user.Auth.Token = nil

	return s.UpdateUser(ctx, user)
}

// SetUserPassword sets a new password for a user, generates a new auth token, and updates the user in the repository.
// It takes a context, a user object, and a password input DTO as parameters.
// Returns an error if password hashing or user update fails.
func (s *userRepository) SetUserPassword(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	user.Auth.Token = new(string)
	user.Auth.Password = new(string)

	hash, err := bcrypt.GenerateFromPassword([]byte(*pass.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	*user.Auth.Token = uuid.New().String()
	*user.Auth.Password = string(hash)

	return s.UpdateUser(ctx, user)
}

// SetUserPhoto sets a user's photo by uploading a file to MinIO and updating the user's photo path in the database.
// If a photo already exists, it deletes the existing photo before uploading the new one.
func (s *userRepository) SetUserPhoto(ctx context.Context, user *domain.User, p *domain.File) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if user.PhotoPath == nil {
			user.PhotoPath = new(string)
		} else if err := s.minio.DeleteObject(ctx, os.Getenv("MINIO_BUCKET_FILES"), *user.PhotoPath); err != nil {
			return err
		}

		*user.PhotoPath = fmt.Sprintf("photos/%v%v", user.ID, p.Extension)
		if err := s.minio.PutObject(ctx, os.Getenv("MINIO_BUCKET_FILES"), *user.PhotoPath, p.File); err != nil {
			return err
		}

		return s.UpdateUser(ctx, user)
	})
}

// GenerateUserPhotoURL generates a URL for the user's photo if a photo path exists. Returns the URL or an error.
func (s *userRepository) GenerateUserPhotoURL(ctx context.Context, user *domain.User) (string, error) {
	if user.PhotoPath == nil {
		return "", myerrors.ErrUserHasNoPhoto
	}

	return s.minio.GenerateObjectURL(
		ctx,
		os.Getenv("MINIO_BUCKET_FILES"),
		*user.PhotoPath,
		fmt.Sprintf("%v%v", strings.ToLower(user.Name), filepath.Ext(*user.PhotoPath)),
	)
}
