package domain

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raulaguila/go-api/internal/pkg/dto"
	"github.com/raulaguila/go-api/pkg/pgutils"
)

func TestProfile_TableName(t *testing.T) {
	profile := new(Profile)
	assert.Equal(t, ProfileTableName, profile.TableName())
}

func TestProfile_ToMap(t *testing.T) {
	tests := []struct {
		name  string
		input *Profile
		want  map[string]any
	}{
		{
			"ValidProfile",
			&Profile{Name: "John Doe", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true, "write": false}`)))},
			map[string]any{"name": "John Doe", "permissions": pgutils.JSONB(json.RawMessage([]byte(`{"read": true, "write": false}`)))},
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
	permissions := pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))
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
			&dto.ProfileInputDTO{Name: new(string), Permissions: &permissions},
			&Profile{},
			&Profile{Name: "New Name", Permissions: pgutils.JSONB(json.RawMessage([]byte(`{"read": true}`)))},
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
