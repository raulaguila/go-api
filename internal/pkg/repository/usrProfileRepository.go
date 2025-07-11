package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
)

func NewProfileRepository(postgreDB *gorm.DB) domain.ProfileRepository {
	return &profileRepository{
		postgreDB: postgreDB,
	}
}

type profileRepository struct {
	postgreDB *gorm.DB
}

func (s *profileRepository) applyFilter(ctx context.Context, f *dto.ProfileFilter) *gorm.DB {
	postgreDB := s.postgreDB.WithContext(ctx)
	if f != nil {
		if f.ID != nil {
			postgreDB = postgreDB.Where("id = ?", *f.ID)
		}

		if where := f.ApplySearchLike("name"); where != "" {
			postgreDB = postgreDB.Where(where)
		}
		postgreDB = postgreDB.Order(f.ApplyOrder(nil))
		if f.WithPermissions != nil && !(*f.WithPermissions) {
			postgreDB = postgreDB.Omit("permissions")
		}
		if !f.ListRoot {
			postgreDB = postgreDB.Where("name != ?", "ROOT")
		}

		postgreDB = postgreDB.Order(f.ApplyOrder(nil))
	}

	return postgreDB.Group("id")
}

func (s *profileRepository) CountProfiles(ctx context.Context, f *dto.ProfileFilter) (int64, error) {
	var count int64
	return count, s.applyFilter(ctx, f).Model(new(domain.Profile)).Count(&count).Error
}

func (s *profileRepository) GetProfiles(ctx context.Context, f *dto.ProfileFilter) (*[]domain.Profile, error) {
	postgreDB := s.applyFilter(ctx, f)
	if f != nil {
		if ok, offset, limit := f.ApplyPagination(); ok {
			postgreDB = postgreDB.Offset(offset).Limit(limit)
		}
	}

	profiles := new([]domain.Profile)
	return profiles, postgreDB.Find(profiles).Error
}

func (s *profileRepository) GetProfile(ctx context.Context, input *domain.Profile) error {
	return s.postgreDB.WithContext(ctx).Where(input).First(input).Error
}

func (s *profileRepository) CreateProfile(ctx context.Context, input *domain.Profile) error {
	return s.postgreDB.WithContext(ctx).Create(input).Error
}

func (s *profileRepository) UpdateProfile(ctx context.Context, input *domain.Profile) error {
	return s.postgreDB.WithContext(ctx).Model(input).Updates(input.ToMap()).Error
}

func (s *profileRepository) DeleteProfiles(ctx context.Context, ids []uint) error {
	return s.postgreDB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(new(domain.Profile), ids)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
