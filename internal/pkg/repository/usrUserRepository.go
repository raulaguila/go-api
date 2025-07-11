package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/packhub"
	"github.com/raulaguila/go-api/pkg/utils"
)

func NewUserRepository(postgreDB *gorm.DB) domain.UserRepository {
	return &userRepository{
		postgreDB: postgreDB,
	}
}

type userRepository struct {
	postgreDB *gorm.DB
}

func (s *userRepository) applyFilter(ctx context.Context, f *dto.UserFilter) *gorm.DB {
	postgreDB := s.postgreDB.WithContext(ctx)
	if f != nil {
		if f.ID != nil {
			postgreDB = postgreDB.Where(domain.UserTableName+".id = ?", *f.ID)
		}

		if f.Status != nil {
			postgreDB = postgreDB.Where(domain.AuthTableName+".status = ?", *f.Status)
		}

		if f.ProfileID != 0 {
			postgreDB = postgreDB.Where(domain.AuthTableName+".profile_id = ?", f.ProfileID)
		}

		postgreDB = postgreDB.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.auth_id", domain.AuthTableName, domain.AuthTableName, domain.UserTableName))
		postgreDB = postgreDB.Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.profile_id", domain.ProfileTableName, domain.ProfileTableName, domain.AuthTableName))

		if where := f.ApplySearchLike(
			domain.UserTableName+".name",
			domain.UserTableName+".username",
			domain.UserTableName+".mail",
			domain.ProfileTableName+".name",
		); where != "" {
			postgreDB = postgreDB.Where(where)
		}

		postgreDB = postgreDB.Order(f.ApplyOrder(packhub.Pointer(domain.UserTableName)))
	}

	return postgreDB.Group(domain.UserTableName + ".id")
}

func (s *userRepository) CountUsers(ctx context.Context, f *dto.UserFilter) (int64, error) {
	var count int64
	return count, s.applyFilter(ctx, f).Model(new(domain.User)).Count(&count).Error
}

func (s *userRepository) GetUsers(ctx context.Context, f *dto.UserFilter) (*[]domain.User, error) {
	postgreDB := s.applyFilter(ctx, f)
	if f != nil {
		if ok, offset, limit := f.ApplyPagination(); ok {
			postgreDB = postgreDB.Offset(offset).Limit(limit)
		}
	}

	users := new([]domain.User)
	return users, postgreDB.Preload(utils.PGAuthProfile).Find(users).Error
}

func (s *userRepository) GetUser(ctx context.Context, input *domain.User) error {
	return s.postgreDB.WithContext(ctx).Where(input).Preload(utils.PGAuthProfile).First(input).Error
}

func (s *userRepository) GetUserByToken(ctx context.Context, token string) (*domain.User, error) {
	user := new(domain.User)
	return user, s.postgreDB.
		WithContext(ctx).
		Joins(fmt.Sprintf("JOIN %v ON %v.id = %v.auth_id", domain.AuthTableName, domain.AuthTableName, domain.UserTableName)).
		Preload(utils.PGAuthProfile).
		First(user, domain.AuthTableName+".token = ?", token).Error
}

func (s *userRepository) CreateUser(ctx context.Context, input *domain.User) error {
	return s.postgreDB.Session(&gorm.Session{FullSaveAssociations: true}).WithContext(ctx).Create(input).Error
}

func (s *userRepository) UpdateUser(ctx context.Context, input *domain.User) error {
	return s.postgreDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := s.postgreDB.Model(input.Auth).Updates(input.Auth.ToMap()).Error; err != nil {
			return err
		}

		return s.postgreDB.Model(input).Updates(input.ToMap()).Error
	})
}

func (s *userRepository) DeleteUsers(ctx context.Context, toDelete []uint) error {
	users := new([]domain.User)
	if err := s.postgreDB.WithContext(ctx).Find(users, toDelete).Error; err != nil {
		return err
	}
	if len(*users) == 0 {
		return gorm.ErrRecordNotFound
	}

	return s.postgreDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Select(clause.Associations).Where(users).Delete(users)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
