package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/filter"
)

func NewProfileRepository(db *gorm.DB) domain.ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

type profileRepository struct {
	db *gorm.DB
}

func (s *profileRepository) applyFilter(ctx context.Context, f *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if f != nil {
		db = f.ApplySearchLike(db, "name")
		db = f.ApplyOrder(db, nil)
	}

	return db
}

func (s *profileRepository) CountProfiles(ctx context.Context, f *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, f)

	return count, db.Model(new(domain.Profile)).Count(&count).Error
}

func (s *profileRepository) GetProfiles(ctx context.Context, f *filter.Filter) (*[]domain.Profile, error) {
	db := s.applyFilter(ctx, f)
	if f != nil {
		db = f.ApplyPagination(db)
	}

	profiles := new([]domain.Profile)
	return profiles, db.Find(profiles).Error
}

func (s *profileRepository) GetProfile(ctx context.Context, p *domain.Profile) error {
	return s.db.WithContext(ctx).Where(p).First(p).Error
}

func (s *profileRepository) CreateProfile(ctx context.Context, p *domain.Profile) error {
	return s.db.WithContext(ctx).Create(p).Error
}

func (s *profileRepository) UpdateProfile(ctx context.Context, p *domain.Profile) error {
	return s.db.WithContext(ctx).Model(p).Updates(p.ToMap()).Error
}

func (s *profileRepository) DeleteProfiles(ctx context.Context, ids []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
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
