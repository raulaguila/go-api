package repository

import (
	"context"

	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/filter"
)

// NewProfileRepository creates a new instance of the ProfileRepository using the provided gorm.DB connection.
// It returns a concrete implementation of the ProfileRepository interface.
func NewProfileRepository(db *gorm.DB) domain.ProfileRepository {
	return &profileRepository{
		db: db,
	}
}

// profileRepository is a struct that provides methods to perform CRUD operations on Profile entities using a gorm.DB instance.
type profileRepository struct {
	db *gorm.DB
}

// applyFilter applies search and ordering filters to the database query using provided filter parameters.
func (s *profileRepository) applyFilter(ctx context.Context, filter *filter.Filter) *gorm.DB {
	db := s.db.WithContext(ctx)
	if filter != nil {
		db = filter.ApplySearchLike(db, "name")
		db = filter.ApplyOrder(db, nil)
	}

	return db
}

// CountProfiles counts the number of profiles in the database applying the specified filter.
func (s *profileRepository) CountProfiles(ctx context.Context, filter *filter.Filter) (int64, error) {
	var count int64
	db := s.applyFilter(ctx, filter)
	return count, db.Model(new(domain.Profile)).Count(&count).Error
}

// GetProfiles retrieves a list of profiles from the database based on provided filter conditions.
// It applies filtering, sorting, and pagination as specified in the filter parameter.
func (s *profileRepository) GetProfiles(ctx context.Context, filter *filter.Filter) (*[]domain.Profile, error) {
	db := s.applyFilter(ctx, filter)
	if filter != nil {
		db = filter.ApplyPagination(db)
	}

	profiles := new([]domain.Profile)
	return profiles, db.Find(profiles).Error
}

// GetProfileByID retrieves a Profile from the database using the provided profileID.
// It takes a context and a uint profileID as parameters and returns a pointer to the Profile and an error if any occurs.
func (s *profileRepository) GetProfileByID(ctx context.Context, profileID uint) (*domain.Profile, error) {
	profile := new(domain.Profile)
	return profile, s.db.WithContext(ctx).First(profile, profileID).Error
}

// CreateProfile saves a new profile into the database using the provided context and profile details.
// If an error occurs during the creation process, it returns the error.
func (s *profileRepository) CreateProfile(ctx context.Context, profile *domain.Profile) error {
	return s.db.WithContext(ctx).Create(profile).Error
}

// UpdateProfile updates an existing profile in the database with the given data.
// It applies the context and updates the profile fields using a map representation.
func (s *profileRepository) UpdateProfile(ctx context.Context, profile *domain.Profile) error {
	return s.db.WithContext(ctx).Model(profile).Updates(profile.ToMap()).Error
}

// DeleteProfiles deletes profiles from the database using the provided list of IDs.
// It executes the delete operation within a transaction to ensure atomicity.
// The function accepts a context for request-scoped values and cancellation.
// The toDelete parameter is a slice of profile IDs that should be removed.
// Returns an error if the operation fails, otherwise returns nil.
func (s *profileRepository) DeleteProfiles(ctx context.Context, toDelete []uint) error {
	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return tx.Delete(new(domain.Profile), toDelete).Error
	})
}
