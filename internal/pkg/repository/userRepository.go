package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

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
	tx := s.db.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Model(user).Updates(user.ToMap())
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
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
