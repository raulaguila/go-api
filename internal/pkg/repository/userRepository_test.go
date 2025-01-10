package repository

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/raulaguila/go-api/internal/pkg/filters"
	"github.com/raulaguila/go-api/pkg/filter"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"gorm.io/gorm/logger"
	"testing"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/pkg/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepository_CountUsers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewUserRepository(db)

	tests := []struct {
		name          string
		mockSetup     func()
		filter        *filters.UserFilter
		expectedCount int64
		expectedError error
	}{
		{
			name: "success_count_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com"}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com"}).Error)
			},
			filter:        nil,
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "success_count_0",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
			},
			filter:        nil,
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			count, err := repository.CountUsers(context.Background(), tt.filter)

			assert.Equal(t, tt.expectedCount, count)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestUserRepository_GetUsers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewUserRepository(db)

	tests := []struct {
		name          string
		mockSetup     func()
		filter        *filters.UserFilter
		expectedNames []string
		expectedErr   error
	}{
		{
			name: "success_get_users_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			filter:        &filters.UserFilter{Filter: *filter.New("name", "asc")},
			expectedNames: []string{"User 1", "User 2"},
			expectedErr:   nil,
		},
		{
			name: "success_get_users_0",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
			},
			filter:        &filters.UserFilter{Filter: *filter.New("name", "asc")},
			expectedNames: []string{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			data, err := repository.GetUsers(context.Background(), tt.filter)

			assert.Equal(t, tt.expectedErr, err)
			for i, name := range tt.expectedNames {
				assert.Equal(t, name, (*data)[i].Name)
			}
		})
	}
}

func TestUserRepository_GetUserByID(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewUserRepository(db)

	tests := []struct {
		name         string
		mockSetup    func()
		userInput    *domain.User
		expectedName string
		expectedErr  error
	}{
		{
			name: "success_get_user_1",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			userInput:    &domain.User{Base: domain.Base{ID: 1}},
			expectedName: "User 1",
			expectedErr:  nil,
		},
		{
			name: "success_get_user_2",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			userInput:    &domain.User{Base: domain.Base{ID: 2}},
			expectedName: "User 2",
			expectedErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repository.GetUser(context.Background(), tt.userInput)

			assert.Equal(t, tt.expectedErr, err)
			assert.Equal(t, tt.expectedName, tt.userInput.Name)
		})
	}
}

func TestUserRepository_CreateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	utils.PanicIfErr(err)
	repository := NewUserRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		input       *domain.User
		expectedErr error
	}{
		{
			name: "success_create_user",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
			},
			input:       &domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}},
			expectedErr: nil,
		},
		{
			name: "error_duplicated_user",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			input:       &domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}},
			expectedErr: errors.New("UNIQUE constraint failed: users.mail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err = repository.CreateUser(context.Background(), tt.input)

			if tt.expectedErr == nil {
				assert.NoError(t, err)
			} else {
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}

func TestUserRepository_UpdateUser(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewUserRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		input       *domain.User
		expectedErr error
	}{
		{
			name: "success_update_user",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			input:       &domain.User{Base: domain.Base{ID: 1}, Name: "Updated User 1", Email: "user1@email.com", Auth: &domain.Auth{Base: domain.Base{ID: 1}, Status: false, ProfileID: 1}},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err := repository.UpdateUser(context.Background(), tt.input)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestUserRepository_DeleteUsers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewUserRepository(db)

	tests := []struct {
		name        string
		mockSetup   func()
		toDelete    []uint
		expectedErr error
	}{
		{
			name: "success_delete_users",
			mockSetup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 3", Email: "user3@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			toDelete:    []uint{1, 2, 3},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			err = repository.DeleteUsers(context.Background(), tt.toDelete)

			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
