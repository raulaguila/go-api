package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/raulaguila/packhub"
	
	"github.com/raulaguila/go-api/internal/pkg/domain"
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

func (s *userRepository) applyFilter(ctx context.Context, f *filters.UserFilter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if f != nil {
		if f.ProfileID != 0 {
			db = db.Where(domain.AuthTableName+".profile_id = ?", f.ProfileID)
		}
		db = db.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.auth_id", domain.AuthTableName, domain.AuthTableName, domain.UserTableName))
		db = db.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.profile_id", domain.ProfileTableName, domain.ProfileTableName, domain.AuthTableName))
		if where := f.ApplySearchLike(
			domain.UserTableName+".name",
			domain.UserTableName+".mail",
			domain.ProfileTableName+".name",
		); where != "" {
			db = db.Where(where)
		}
		db = db.Order(f.ApplyOrder(packhub.Pointer(domain.UserTableName)))
	}

	return db.Group(domain.UserTableName + ".id")
}

func (s *userRepository) CountUsers(ctx context.Context, f *filters.UserFilter) (int64, error) {
	var count int64
	return count, s.applyFilter(ctx, f).Model(new(domain.User)).Count(&count).Error
}

func (s *userRepository) GetUsers(ctx context.Context, f *filters.UserFilter) (*[]domain.User, error) {
	db := s.applyFilter(ctx, f)
	if f != nil {
		if ok, offset, limit := f.ApplyPagination(); ok {
			db = db.Offset(offset).Limit(limit)
		}
	}

	users := new([]domain.User)
	return users, db.Preload(utils.PGAuthProfile).Find(users).Error
}

func (s *userRepository) GetUser(ctx context.Context, user *domain.User) error {
	return s.db.WithContext(ctx).Where(user).Preload(utils.PGAuthProfile).First(user).Error
}

func (s *userRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	user := new(domain.User)
	return user, s.db.
		WithContext(ctx).
		Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.auth_id", domain.AuthTableName, domain.AuthTableName, domain.UserTableName)).
		Preload(utils.PGAuthProfile).
		First(user, domain.AuthTableName+".token = ?", token).Error
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
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(new(domain.User), "id IN ?", toDelete)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
