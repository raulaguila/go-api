package service

import (
	"context"
	"github.com/raulaguila/go-api/internal/pkg/mocks"
	"gorm.io/gorm"
	"testing"

	"github.com/raulaguila/go-api/internal/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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
					On("GetUserByID", mock.Anything, uint(1)).
					Return(&domain.User{
						Base: domain.Base{
							ID: 1,
						},
						Auth: &domain.Auth{
							Status: false,
							Profile: &domain.Profile{
								Base: domain.Base{
									ID: 1,
								},
								Name: "ADMIN",
							},
						},
						Name:  "John Doe",
						Email: "johndoe@example.com",
					}, nil).
					Once()
			},
			userID:  1,
			wantErr: false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetUserByID", mock.Anything, uint(1)).
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
			_, err := service.GetUserByID(context.Background(), tt.userID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

//func TestUserService_CreateUser(t *testing.T) {
//	mockRepository := new(mocks.UserRepositoryMock)
//	service := NewUserService(mockRepository)
//
//	tests := []struct {
//		name      string
//		setupMock func()
//		userInput *dto.UserInputDTO
//		wantErr   bool
//	}{
//		{
//			name: "success",
//			setupMock: func() {
//				mockRepository.
//					On("CreateUser", mock.Anything, mock.Anything).
//					Return(nil).
//					Once()
//				mockRepository.
//					On("GetUserByID", mock.Anything, uint(1)).
//					Return(&domain.User{Base: domain.Base{ID: 1}, Name: "John Doe", Email: "johndoe@example.com"}, nil).
//					Once()
//			},
//			userInput: &dto.UserInputDTO{Name: helper.StringPtr("John Doe"), Email: helper.StringPtr("johndoe@example.com")},
//			wantErr:   false,
//		},
//		{
//			name: "create error",
//			setupMock: func() {
//				mockRepository.
//					On("CreateUser", mock.Anything, mock.Anything).
//					Return(errors.New("create error")).
//					Once()
//			},
//			userInput: &dto.UserInputDTO{Name: helper.StringPtr("John Doe"), Email: helper.StringPtr("johndoe@example.com")},
//			wantErr:   true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMock()
//			_, err := service.CreateUser(context.Background(), tt.userInput)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}

//func TestUserService_SetUserPassword(t *testing.T) {
//	mockRepo := new(mockUserRepository)
//	service := &userService{userRepository: mockRepo}
//
//	tests := []struct {
//		name      string
//		setupMock func()
//		email     string
//		password  *dto.PasswordInputDTO
//		wantErr   bool
//	}{
//		{
//			name: "successful password set",
//			setupMock: func() {
//				mockRepo.On("GetUserByMail", mock.Anything, "test@example.com").
//					Return(&domain.User{Auth: &domain.Auth{Password: nil}}, nil)
//				mockRepo.On("SetUserPassword", mock.Anything, mock.Anything, mock.Anything).Return(nil)
//			},
//			email:    "test@example.com",
//			password: &dto.PasswordInputDTO{Password: "newpassword"},
//			wantErr:  false,
//		},
//		{
//			name: "password already set",
//			setupMock: func() {
//				mockRepo.On("GetUserByMail", mock.Anything, "test@example.com").
//					Return(&domain.User{Auth: &domain.Auth{Password: new(string)}}, nil)
//			},
//			email:    "test@example.com",
//			password: &dto.PasswordInputDTO{Password: "newpassword"},
//			wantErr:  true,
//		},
//		{
//			name: "user not found",
//			setupMock: func() {
//				mockRepo.On("GetUserByMail", mock.Anything, "test@example.com").
//					Return(nil, errors.New("user not found"))
//			},
//			email:    "test@example.com",
//			password: &dto.PasswordInputDTO{Password: "newpassword"},
//			wantErr:  true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMock()
//			err := service.SetUserPassword(context.Background(), tt.email, tt.password)
//			if tt.wantErr {
//				assert.Error(t, err)
//				if tt.name == "password already set" {
//					assert.True(t, errors.Is(err, myerrors.ErrUserHasPass))
//				}
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserService_ResetUserPassword(t *testing.T) {
//	mockRepo := new(mockUserRepository)
//	service := &userService{userRepository: mockRepo}
//
//	tests := []struct {
//		name      string
//		setupMock func()
//		email     string
//		wantErr   bool
//	}{
//		{
//			name: "successful reset",
//			setupMock: func() {
//				mockRepo.On("GetUserByMail", mock.Anything, "test@example.com").
//					Return(&domain.User{Auth: &domain.Auth{Password: new(string)}}, nil)
//				mockRepo.On("ResetUserPassword", mock.Anything, mock.Anything).Return(nil)
//			},
//			email:   "test@example.com",
//			wantErr: false,
//		},
//		{
//			name: "password nil, no reset",
//			setupMock: func() {
//				mockRepo.On("GetUserByMail", mock.Anything, "test@example.com").
//					Return(&domain.User{Auth: &domain.Auth{Password: nil}}, nil)
//			},
//			email:   "test@example.com",
//			wantErr: false,
//		},
//		{
//			name: "user not found",
//			setupMock: func() {
//				mockRepo.On("GetUserByMail", mock.Anything, "test@example.com").
//					Return(nil, errors.New("user not found"))
//			},
//			email:   "test@example.com",
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMock()
//			err := service.ResetUserPassword(context.Background(), tt.email)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserService_DeleteUsers(t *testing.T) {
//	mockRepo := new(mockUserRepository)
//	service := &userService{userRepository: mockRepo}
//
//	tests := []struct {
//		name      string
//		setupMock func()
//		ids       []uint
//		wantErr   bool
//	}{
//		{
//			name: "successful delete",
//			setupMock: func() {
//				mockRepo.On("DeleteUsers", mock.Anything, []uint{1, 2, 3}).
//					Return(nil)
//			},
//			ids:     []uint{1, 2, 3},
//			wantErr: false,
//		},
//		{
//			name: "delete error",
//			setupMock: func() {
//				mockRepo.On("DeleteUsers", mock.Anything, []uint{1, 2, 3}).
//					Return(errors.New("delete error"))
//			},
//			ids:     []uint{1, 2, 3},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMock()
//			err := service.DeleteUsers(context.Background(), tt.ids)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserService_SetUserPhoto(t *testing.T) {
//	mockRepo := new(mockUserRepository)
//	service := &userService{userRepository: mockRepo}
//
//	tests := []struct {
//		name      string
//		setupMock func()
//		userID    uint
//		file      *domain.File
//		wantErr   bool
//	}{
//		{
//			name: "successful set photo",
//			setupMock: func() {
//				mockRepo.On("GetUserByID", mock.Anything, uint(1)).
//					Return(&domain.User{ID: 1}, nil)
//				mockRepo.On("SetUserPhoto", mock.Anything, mock.Anything, mock.Anything).
//					Return(nil)
//			},
//			userID:  1,
//			file:    &domain.File{},
//			wantErr: false,
//		},
//		{
//			name: "user not found",
//			setupMock: func() {
//				mockRepo.On("GetUserByID", mock.Anything, uint(1)).
//					Return(nil, errors.New("user not found"))
//			},
//			userID:  1,
//			file:    &domain.File{},
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMock()
//			err := service.SetUserPhoto(context.Background(), tt.userID, tt.file)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
//
//func TestUserService_GenerateUserPhotoURL(t *testing.T) {
//	mockRepo := new(mockUserRepository)
//	service := &userService{userRepository: mockRepo}
//
//	tests := []struct {
//		name      string
//		setupMock func()
//		userID    uint
//		wantErr   bool
//	}{
//		{
//			name: "successful URL generation",
//			setupMock: func() {
//				mockRepo.On("GetUserByID", mock.Anything, uint(1)).
//					Return(&domain.User{ID: 1}, nil)
//				mockRepo.On("GenerateUserPhotoURL", mock.Anything, mock.Anything).
//					Return("http://photo.url", nil)
//			},
//			userID:  1,
//			wantErr: false,
//		},
//		{
//			name: "user not found",
//			setupMock: func() {
//				mockRepo.On("GetUserByID", mock.Anything, uint(1)).
//					Return(nil, errors.New("user not found"))
//			},
//			userID:  1,
//			wantErr: true,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			tt.setupMock()
//			_, err := service.GenerateUserPhotoURL(context.Background(), tt.userID)
//			if tt.wantErr {
//				assert.Error(t, err)
//			} else {
//				assert.NoError(t, err)
//			}
//		})
//	}
//}
