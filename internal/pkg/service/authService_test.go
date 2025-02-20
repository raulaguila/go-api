package service

import (
	"context"
	"os"
	"testing"

	"github.com/lib/pq"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"github.com/raulaguila/go-api/internal/pkg/_mocks"
	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Login(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewAuthService(mockRepository)

	tests := []struct {
		name      string
		setupMock func()
		authInput *dto.AuthInputDTO
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
			authInput: &dto.AuthInputDTO{Login: "admin@admin.com", Password: "12345678"},
			wantErr:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetUser", mock.Anything, mock.Anything).
					Return(gorm.ErrRecordNotFound).
					Once()
			},
			authInput: &dto.AuthInputDTO{Login: "admin@admin.com", Password: "12345678"},
			wantErr:   true,
		},
	}

	_ = os.Setenv("userRepositoryMockType", "authTests")
	defer func() {
		_ = os.Unsetenv("userRepositoryMockType")
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.Login(context.Background(), tt.authInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthService_Me(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewAuthService(mockRepository)

	tests := []struct {
		name      string
		userInput *domain.User
	}{
		{
			name: "success",
			userInput: &domain.User{
				Base:  domain.Base{ID: 1},
				Name:  "Jhon Cena",
				Email: "jhoncena@gmail.com",
				Auth: &domain.Auth{
					Status: false,
					Profile: &domain.Profile{
						Base:        domain.Base{ID: 1},
						Name:        "ADMIN",
						Permissions: pq.StringArray{"read"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsType(t, &dto.UserOutputDTO{}, service.Me(tt.userInput))
		})
	}
}

func TestAuthService_Refresh(t *testing.T) {
	mockRepository := new(_mocks.UserRepositoryMock)
	service := NewAuthService(mockRepository)

	tests := []struct {
		name      string
		userInput *domain.User
	}{
		{
			name: "success",
			userInput: &domain.User{
				Base:  domain.Base{ID: 1},
				Name:  "Jhon Cena",
				Email: "jhoncena@gmail.com",
				Auth: &domain.Auth{
					Status: false,
					Profile: &domain.Profile{
						Base:        domain.Base{ID: 1},
						Name:        "ADMIN",
						Permissions: pq.StringArray{"read"},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsType(t, &dto.AuthOutputDTO{}, service.Refresh(tt.userInput))
		})
	}
}
