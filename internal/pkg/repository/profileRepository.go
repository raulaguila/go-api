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

func (s *profileRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if filter != nil {
		db = filter.ApplySearchLike(db, "name")
		db = filter.ApplyOrder(db, nil)
	}

	return db
}

func (s *profileRepository) CountProfiles(ctx context.Context, filter *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	
	return count, db.Model(new(domain.Profile)).Count(&count).Error
}

func (s *profileRepository) GetProfiles(ctx context.Context, filter *filter.Filter) (*[]domain.Profile, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	profiles := new([]domain.Profile)
	return profiles, db.Find(profiles).Error
}

func (s *profileRepository) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	profile := new(domain.Profile)
	return profile, s.db.WithContext(ctx).First(profile, profileID).Error
}

func (s *profileRepository) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	return s.db.WithContext(ctx).Create(profile).Error
}

func (s *profileRepository) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	return s.db.WithContext(ctx).Model(profile).Updates(profile.ToMap()).Error
}

func (s *profileRepository) DeleteProfiles(ctx context.Context, toDelete []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(new(domain.Profile), toDelete).Error
	})
}
