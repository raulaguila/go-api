package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/utils"
)

func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{
		db: db,
	}
}

type userRepository struct {
	db *gorm.DB
}

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

func (s *userRepository) CountUsers(ctx context.Context, filter *filters.UserFilter) (int64, error) {
	var count int64

	db := s.applyFilter(ctx, filter)
	return count, db.Model(new(domain.User)).Count(&count).Error
}

func (s *userRepository) GetUsers(ctx context.Context, filter *filters.UserFilter) (*[]domain.User, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	users := new([]domain.User)
	return users, db.Preload(utils.PGAuthProfile).Find(users).Error
}

func (s *userRepository) GetUserByID(ctx context.Context, userID uint) (*domain.User, error) {
	user := new(domain.User)
	return user, s.db.WithContext(ctx).Preload(utils.PGAuthProfile).First(user, userID).Error
}

func (s *userRepository) GetUserByMail(ctx context.Context, mail string) (*domain.User, error) {
	user := new(domain.User)
	return user, s.db.WithContext(ctx).Preload(utils.PGAuthProfile).First(user, "mail = ?", mail).Error
}

func (s *userRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	auth := new(domain.Auth)
	if err := s.db.WithContext(ctx).First(auth, "token = ?", token).Error; err != nil {
		return nil, err
	}

	user := new(domain.User)
	return user, s.db.WithContext(ctx).Preload(utils.PGAuthProfile).First(user, "auth_id = ?", auth.ID).Error
}

func (s *userRepository) CreateUser(ctx context.Context, user *domain.User) error {
	return s.db.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).Create(user).Error
}

func (s *userRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.db.Model(user.Auth).Updates(user.Auth.ToMap()).Error; err != nil {
			return err
		}

		return s.db.Model(user).Updates(user.ToMap()).Error
	})
}

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

		return tx.WithContext(ctx).Delete(new(domain.Auth), auths).Error
	})
}

func (s *userRepository) ResetUserPassword(ctx context.Context, user *domain.User) error {
	user.Auth.Password = nil
	user.Auth.Token = nil

	return s.UpdateUser(ctx, user)
}

func (s *userRepository) SetUserPassword(ctx context.Context, user *domain.User, pass *dto.PasswordInputDTO) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(*pass.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Auth.Token = utils.Pointer(uuid.New().String())
	user.Auth.Password = utils.Pointer(string(hash))

	return s.UpdateUser(ctx, user)
}
