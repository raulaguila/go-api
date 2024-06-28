package repository

import (
	"context"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/filter"
	"gorm.io/gorm"
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
	return count, db.Model(&domain.Profile{}).Count(&count).Error
}

func (s *profileRepository) GetProfiles(ctx context.Context, filter *filter.Filter) (*[]domain.Profile, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	profiles := &[]domain.Profile{}
	return profiles, db.Find(profiles).Error
}

func (s *profileRepository) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	profile := &domain.Profile{}
	return profile, s.db.WithContext(ctx).First(profile, profileID).Error
}

func (s *profileRepository) CreateProfile(ctx context.Context, data *dto.ProfileInputDTO) (*domain.Profile, error) {
	profile := &domain.Profile{Permissions: map[string]any{}}
	if err := profile.Bind(data); err != nil {
		return nil, err
	}

	return profile, s.db.WithContext(ctx).Create(profile).Error
}

func (s *profileRepository) UpdateProfile(ctx context.Context, profile *domain.Profile, data *dto.ProfileInputDTO) error {
	if err := profile.Bind(data); err != nil {
		return err
	}

	return s.db.WithContext(ctx).Model(profile).Updates(profile.ToMap()).Error
}

func (s *profileRepository) DeleteProfiles(ctx context.Context, toDelete []uint) error {
	if len(toDelete) == 0 {
		return nil
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(&domain.Profile{}, toDelete).Error
	})
}
