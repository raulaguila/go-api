package domain

import (
	"context"
	"reflect"
	"testing"

	"github.com/raulaguila/go-api/internal/pkg/dto"
)

func TestAuth_TableName(t *testing.T) {
	auth := Auth{}
	expected := AuthTableName

	if result := auth.TableName(); result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestAuth_ToMap(t *testing.T) {
	tests := []struct {
		name     string
		auth     Auth
		expected map[string]any
	}{
		{
			name: "TokenAndPasswordNil",
			auth: Auth{
				Status:    true,
				ProfileID: 1,
				Token:     nil,
				Password:  nil,
			},
			expected: map[string]any{
				"status":     true,
				"profile_id": uint(1),
				"token":      nil,
				"password":   nil,
			},
		},
		{
			name: "TokenAndPasswordSet",
			auth: Auth{
				Status:    false,
				ProfileID: 2,
				Token:     stringPointer("token123"),
				Password:  stringPointer("password123"),
			},
			expected: map[string]any{
				"status":     false,
				"profile_id": uint(2),
				"token":      "token123",
				"password":   "password123",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.auth.ToMap()
			if !reflect.DeepEqual(*result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, *result)
			}
		})
	}
}

func stringPointer(s string) *string {
	return &s
}

type mockAuthService struct{}

func (m *mockAuthService) Login(ctx context.Context, input *dto.AuthInputDTO, ip string) (*dto.AuthOutputDTO, error) {
	// Mock implementation for testing
	return &dto.AuthOutputDTO{}, nil
}

func (m *mockAuthService) Refresh(user *User, ip string) *dto.AuthOutputDTO {
	// Mock implementation for testing
	return &dto.AuthOutputDTO{}
}

func (m *mockAuthService) Me(user *User) *dto.UserOutputDTO {
	// Mock implementation for testing
	return &dto.UserOutputDTO{}
}
