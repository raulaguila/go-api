package service

import (
	"context"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/raulaguila/packhub"

	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
)

func TestUserService_GetUsers(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewUserService(mockRepository)

	tests := []struct {
		name      string
		setupMock func()
		userID    uint
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("GetUsers", mock.Anything, mock.Anything).
					Return(&[]domain.User{{
						Base:  domain.Base{ID: 1},
						Name:  "User 01",
						Email: "user02@example.com",
						Auth: &domain.Auth{
							Status: false,
							Profile: &domain.Profile{
								Base:        domain.Base{ID: 1},
								Name:        "Profile 01",
								Permissions: make(pq.StringArray, 0),
							},
						},
					}}, nil).
					Once()
				mockRepository.
					On("CountUsers", mock.Anything, mock.Anything).
					Return(int64(1), nil).
					Once()
			},
			userID:  1,
			wantErr: false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetUsers", mock.Anything, mock.Anything).
					Return(nil, gorm.ErrRecordNotFound).
					Once()
			},
			userID:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.GetUsers(context.Background(), &dto.UserFilter{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewUserService(mockRepository)

	tests := []struct {
		name      string
		setupMock func()
		userInput *dto.UserInputDTO
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("CreateUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			userInput: &dto.UserInputDTO{Name: packhub.Pointer("John Doe"), Email: packhub.Pointer("johndoe@example.com"), ProfileID: packhub.Pointer(uint(1))},
			wantErr:   false,
		},
		{
			name: "create error",
			setupMock: func() {
				mockRepository.
					On("CreateUser", mock.Anything, mock.Anything).
					Return(gorm.ErrDuplicatedKey).
					Once()
			},
			userInput: &dto.UserInputDTO{Name: packhub.Pointer("John Doe"), Email: packhub.Pointer("johndoe@example.com"), ProfileID: packhub.Pointer(uint(1))},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.CreateUser(context.Background(), tt.userInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewUserService(mockRepository)

	tests := []struct {
		name      string
		setup     func()
		userID    uint
		userInput *dto.UserInputDTO
		wantErr   bool
	}{
		{
			name: "success",
			setup: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(nil).
					Twice()
				mockRepository.
					On("UpdateUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			userID:    1,
			userInput: &dto.UserInputDTO{Name: packhub.Pointer("John Doe"), Email: packhub.Pointer("johndoe@example.com"), ProfileID: packhub.Pointer(uint(1))},
			wantErr:   false,
		},
		{
			name: "create error",
			setup: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("UpdateUser", mock.Anything, mock.Anything).
					Return(gorm.ErrDuplicatedKey).
					Once()
			},
			userID:    1,
			userInput: &dto.UserInputDTO{Name: packhub.Pointer("John Doe"), Email: packhub.Pointer("johndoe@example.com"), ProfileID: packhub.Pointer(uint(1))},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := service.UpdateUser(context.Background(), tt.userID, tt.userInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_DeleteUsers(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewUserService(mockRepository)

	tests := []struct {
		name     string
		setup    func()
		idsInput []uint
		wantErr  bool
	}{
		{
			name: "success",
			setup: func() {
				mockRepository.
					On("DeleteUsers", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			idsInput: []uint{1},
			wantErr:  false,
		},
		{
			name: "not found",
			setup: func() {
				mockRepository.
					On("DeleteUsers", mock.Anything, mock.Anything).
					Return(gorm.ErrRecordNotFound).
					Once()
			},
			idsInput: []uint{1},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.DeleteUsers(context.Background(), tt.idsInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_ResetUserPassword(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewUserService(mockRepository)

	tests := []struct {
		name      string
		setup     func()
		mailInput string
		wantErr   bool
	}{
		{
			name: "success",
			setup: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("UpdateUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			mailInput: "johndoe@example.com",
			wantErr:   false,
		},
		{
			name: "not found",
			setup: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(gorm.ErrRecordNotFound).
					Once()
			},
			mailInput: "johndoe@example.com",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.ResetUserPassword(context.Background(), tt.mailInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProductService_SetUserPassword(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewUserService(mockRepository)

	tests := []struct {
		name      string
		setup     func()
		mailInput string
		passInput *dto.PasswordInputDTO
		wantErr   bool
	}{
		{
			name: "success",
			setup: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("UpdateUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			mailInput: "johndoe@example.com",
			passInput: &dto.PasswordInputDTO{Password: packhub.Pointer("secret"), PasswordConfirm: packhub.Pointer("secret")},
			wantErr:   false,
		},
		{
			name: "not found",
			setup: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(gorm.ErrRecordNotFound).
					Once()
			},
			mailInput: "johndoe@example.com",
			passInput: &dto.PasswordInputDTO{Password: packhub.Pointer("secret"), PasswordConfirm: packhub.Pointer("secret")},
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			err := service.SetUserPassword(context.Background(), tt.mailInput, tt.passInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
