package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/raulaguila/packhub"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
)

func TestUserService_GetUserByID(t *testing.T) {
	mockRepository := new(mocks.UserRepositoryMock)
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
					On("GetUser", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			userID:  1,
			wantErr: false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(gorm.ErrRecordNotFound).
					Once()
			},
			userID:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.GetUserByID(context.Background(), tt.userID)
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}

func TestUserService_CreateUser(t *testing.T) {
	mockRepository := new(mocks.UserRepositoryMock)
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
			assert.Equal(t, err != nil, tt.wantErr)
		})
	}
}
