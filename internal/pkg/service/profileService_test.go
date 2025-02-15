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
	"github.com/raulaguila/go-api/pkg/pgfilter"
)

func TestProfileService_GetProfileByID(t *testing.T) {
	mockRepository := new(_mocks.ProfileRepositoryMock)
	service := NewProfileService(mockRepository)

	tests := []struct {
		name      string
		setupMock func()
		profileID uint
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("GetProfile", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			profileID: 1,
			wantErr:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetProfile", mock.Anything, mock.Anything).
					Return(gorm.ErrRecordNotFound).
					Once()
			},
			profileID: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.GetProfileByID(context.Background(), tt.profileID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProfileService_GetProfiles(t *testing.T) {
	mockRepository := new(_mocks.ProfileRepositoryMock)
	service := NewProfileService(mockRepository)

	tests := []struct {
		name      string
		setupMock func()
		profileID uint
		wantErr   bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("GetProfiles", mock.Anything, mock.Anything).
					Return(&[]domain.Profile{{
						Base:        domain.Base{ID: 1},
						Name:        "Profile 01",
						Permissions: pq.StringArray{"read"},
					}}, nil).
					Once()
				mockRepository.
					On("CountProfiles", mock.Anything, mock.Anything).
					Return(int64(1), nil).
					Once()
			},
			profileID: 1,
			wantErr:   false,
		},
		{
			name: "not found",
			setupMock: func() {
				mockRepository.
					On("GetProfiles", mock.Anything, mock.Anything).
					Return(nil, gorm.ErrRecordNotFound).
					Once()
			},
			profileID: 1,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.GetProfiles(context.Background(), &pgfilter.Filter{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProfileService_CreateProfile(t *testing.T) {
	mockRepository := new(_mocks.ProfileRepositoryMock)
	service := NewProfileService(mockRepository)

	tests := []struct {
		name         string
		setupMock    func()
		profileInput *dto.ProfileInputDTO
		wantErr      bool
	}{
		{
			name: "success",
			setupMock: func() {
				mockRepository.
					On("CreateProfile", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("GetProfile", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			profileInput: &dto.ProfileInputDTO{Name: packhub.Pointer("John Doe"), Permissions: &pq.StringArray{"read"}},
			wantErr:      false,
		},
		{
			name: "create error",
			setupMock: func() {
				mockRepository.
					On("CreateProfile", mock.Anything, mock.Anything).
					Return(gorm.ErrDuplicatedKey).
					Once()
			},
			profileInput: &dto.ProfileInputDTO{Name: packhub.Pointer("John Doe"), Permissions: &pq.StringArray{"read"}},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock()
			_, err := service.CreateProfile(context.Background(), tt.profileInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProfileService_UpdateProfile(t *testing.T) {
	mockRepository := new(_mocks.ProfileRepositoryMock)
	service := NewProfileService(mockRepository)

	tests := []struct {
		name         string
		setup        func()
		profileID    uint
		profileInput *dto.ProfileInputDTO
		wantErr      bool
	}{
		{
			name: "success",
			setup: func() {
				mockRepository.
					On("GetProfile", mock.Anything, mock.Anything).
					Return(nil).
					Twice()
				mockRepository.
					On("UpdateProfile", mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			profileID:    1,
			profileInput: &dto.ProfileInputDTO{Name: packhub.Pointer("John Doe"), Permissions: &pq.StringArray{"read"}},
			wantErr:      false,
		},
		{
			name: "create error",
			setup: func() {
				mockRepository.
					On("GetProfile", mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepository.
					On("UpdateProfile", mock.Anything, mock.Anything).
					Return(gorm.ErrDuplicatedKey).
					Once()
			},
			profileID:    1,
			profileInput: &dto.ProfileInputDTO{Name: packhub.Pointer("John Doe"), Permissions: &pq.StringArray{"read"}},
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			_, err := service.UpdateProfile(context.Background(), tt.profileID, tt.profileInput)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
