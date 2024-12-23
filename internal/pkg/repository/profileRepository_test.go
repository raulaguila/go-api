package repository

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/raulaguila/go-api/pkg/utils"
)

func TestProfileRepository_CountProfiles(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProfileRepository(db)

	tests := []struct {
		name          string
		mockSetup     func()
		filter        *filter.Filter
		expectedCount int64
		expectedError error
	}{
		{
			name: "success_count_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 2", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			filter:        nil,
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "success_count_0",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
			},
			filter:        nil,
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			count, err := repository.CountProfiles(context.Background(), tt.filter)

			assert.Equal(t, tt.expectedCount, count)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestProfileRepository_GetProfiles(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProfileRepository(db)

	tests := []struct {
		name          string
		mockSetup     func()
		filter        *filter.Filter
		expectedNames []string
		expectedErr   error
	}{
		{
			name: "success_get_profiles_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 2", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			filter:        filter.New("name", "asc"),
			expectedNames: []string{"Profile 1", "Profile 2"},
			expectedErr:   nil,
		},
		{
			name: "success_get_profiles_0",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
			},
			filter:        filter.New("name", "asc"),
			expectedNames: []string{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			data, err := repository.GetProfiles(context.Background(), tt.filter)

			for i, name := range tt.expectedNames {
				assert.Equal(t, name, (*data)[i].Name)
			}
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestProfileRepository_GetProfileByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProfileRepository(db)

	tests := []struct {
		name         string
		mockSetup    func()
		profileID    uint
		expectedName string
		expectedErr  error
	}{
		{
			name: "success_get_profile_1",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 2", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			profileID:    1,
			expectedName: "Profile 1",
			expectedErr:  nil,
		},
		{
			name: "success_get_profile_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 2", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			profileID:    2,
			expectedName: "Profile 2",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			data, err := repository.GetProfileByID(context.Background(), tt.profileID)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedName, data.Name)
		})
	}
}

func TestProfileRepository_CreateProfile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	utils.PanicIfErr(err)
	repository := NewProfileRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		input       *domain.Profile
		expectedErr error
	}{
		{
			name: "success_create_profile",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
			},
			input:       &domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))},
			expectedErr: nil,
		},
		{
			name: "error_duplicated_profile",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			input:       &domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))},
			expectedErr: errors.New("UNIQUE constraint failed: users_profile.name"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err = repository.CreateProfile(context.Background(), tt.input)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestProfileRepository_UpdateProfile(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProfileRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		input       *domain.Profile
		expectedErr error
	}{
		{
			name: "success_update_profile",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			input:       &domain.Profile{Base: domain.Base{ID: 1}, Name: "Updated Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repository.UpdateProfile(context.Background(), tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestProfileRepository_DeleteProfiles(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewProfileRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		toDelete    []uint
		expectedErr error
	}{
		{
			name: "success_delete_profiles",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 2", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 3", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			toDelete:    []uint{1, 2, 3},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err = repository.DeleteProfiles(context.Background(), tt.toDelete)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
