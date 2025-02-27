package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgfilter"
	"github.com/raulaguila/go-api/pkg/utils"
)

func TestUserRepository_CountUsers(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	utils.PanicIfErr(err)
	repository := NewUserRepository(db)

	tests := []struct {
		name          string
		setup         func()
		filter        *dto.UserFilter
		expectedCount int64
		expectedError error
	}{
		{
			name: "success_count_2",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com"}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com"}).Error)
			},
			filter:        nil,
			expectedCount: 2,
			expectedError: nil,
		},
		{
			name: "success_count_0",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
			},
			filter:        nil,
			expectedCount: 0,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
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
		setup         func()
		filter        *dto.UserFilter
		expectedNames []string
		expectedErr   error
	}{
		{
			name: "success_get_users_2",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			filter:        &dto.UserFilter{Filter: *pgfilter.New("name", "asc")},
			expectedNames: []string{"User 1", "User 2"},
			expectedErr:   nil,
		},
		{
			name: "success_get_users_0",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
			},
			filter:        &dto.UserFilter{Filter: *pgfilter.New("name", "asc")},
			expectedNames: []string{},
			expectedErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
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
		setup        func()
		userInput    *domain.User
		expectedName string
		expectedErr  error
	}{
		{
			name: "success_get_user_1",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 2", Email: "user2@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			userInput:    &domain.User{Base: domain.Base{ID: 1}},
			expectedName: "User 1",
			expectedErr:  nil,
		},
		{
			name: "success_get_user_2",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pq.StringArray{"read"}}).Error)
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
			tt.setup()
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
		setup       func()
		input       *domain.User
		expectedErr error
	}{
		{
			name: "success_create_user",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pq.StringArray{"read"}}).Error)
			},
			input:       &domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}},
			expectedErr: nil,
		},
		{
			name: "error_duplicated_user",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			input:       &domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}},
			expectedErr: errors.New("UNIQUE constraint failed: users.mail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
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
		setup       func()
		input       *domain.User
		expectedErr error
	}{
		{
			name: "success_update_user",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Auth{}))
				utils.PanicIfErr(db.Migrator().DropTable(&domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.User{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Auth{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", Auth: &domain.Auth{Status: false, ProfileID: 1}}).Error)
			},
			input:       &domain.User{Base: domain.Base{ID: 1}, Name: "Updated User 1", Email: "user1@email.com", Auth: &domain.Auth{Base: domain.Base{ID: 1}, Status: false, ProfileID: 1}},
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
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
		name     string
		setup    func()
		toDelete []uint
	}{
		{
			name: "success_delete_users",
			setup: func() {
				utils.PanicIfErr(db.Migrator().DropTable(&domain.User{}, &domain.Auth{}, &domain.Profile{}))
				utils.PanicIfErr(db.AutoMigrate(&domain.Profile{}, &domain.Auth{}, &domain.User{}))
				utils.PanicIfErr(db.Create(&domain.Profile{Name: "Profile 1", Permissions: pq.StringArray{"read"}}).Error)
				utils.PanicIfErr(db.Create(&domain.Auth{Status: false, ProfileID: 1}).Error)
				utils.PanicIfErr(db.Create(&domain.User{Name: "User 1", Email: "user1@email.com", AuthID: 1}).Error)
			},
			toDelete: []uint{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err = repository.DeleteUsers(context.Background(), tt.toDelete)

			assert.NoError(t, err)
		})
	}
}
