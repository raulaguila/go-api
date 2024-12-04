package domain

import (
	"context"
	"github.com/raulaguila/go-api/pkg/filter"
	"testing"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockProfileRepository is a mock implementation of ProfileRepository
type MockProfileRepository struct {
	mock.Mock
}

func (m *MockProfileRepository) CountProfiles(ctx context.Context, f *filter.Filter) (int64, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockProfileRepository) GetProfileByID(ctx context.Context, id uint) (*Profile, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Profile), args.Error(1)
}

func (m *MockProfileRepository) GetProfiles(ctx context.Context, f *filter.Filter) (*[]Profile, error) {
	args := m.Called(ctx, f)
	return args.Get(0).(*[]Profile), args.Error(1)
}

func (m *MockProfileRepository) CreateProfile(ctx context.Context, p *Profile) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProfileRepository) UpdateProfile(ctx context.Context, p *Profile) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockProfileRepository) DeleteProfiles(ctx context.Context, ids []uint) error {
	args := m.Called(ctx, ids)
	return args.Error(0)
}

func TestProfile_TableName(t *testing.T) {
	profile := &Profile{}
	expected := ProfileTableName
	actual := profile.TableName()
	assert.Equal(t, expected, actual)
}

func TestProfile_ToMap(t *testing.T) {
	tests := []struct {
		name  string
		input *Profile
		want  map[string]any
	}{
		{
			"ValidProfile",
			&Profile{Name: "John Doe", Permissions: pgutils.JSONB{"read": true, "write": false}},
			map[string]any{"name": "John Doe", "permissions": pgutils.JSONB{"read": true, "write": false}},
		},
		{
			"EmptyProfile",
			&Profile{},
			map[string]any{"name": "", "permissions": pgutils.JSONB(nil)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.input.ToMap()
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestProfile_Bind(t *testing.T) {
	tests := []struct {
		name        string
		profileName string
		inputDTO    *dto.ProfileInputDTO
		startData   *Profile
		expected    *Profile
		wantErr     bool
	}{
		{
			"ValidInputDTO",
			"New Name",
			&dto.ProfileInputDTO{Name: new(string), Permissions: pgutils.JSONB{"create": true}},
			&Profile{},
			&Profile{Name: "New Name", Permissions: pgutils.JSONB{"create": true}},
			false,
		},
		{
			"EmptyInputDTO",
			"",
			&dto.ProfileInputDTO{Name: new(string)},
			&Profile{},
			&Profile{},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			*tt.inputDTO.Name = tt.profileName
			err := tt.startData.Bind(tt.inputDTO)
			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, tt.expected, tt.startData)
		})
	}
}
