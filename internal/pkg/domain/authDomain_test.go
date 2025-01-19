package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/raulaguila/packhub"
)

func TestAuth_TableName(t *testing.T) {
	auth := new(Auth)
	assert.Equal(t, AuthTableName, auth.TableName())
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
				"token":      (*string)(nil),
				"password":   (*string)(nil),
			},
		},
		{
			name: "TokenAndPasswordSet",
			auth: Auth{
				Status:    false,
				ProfileID: 2,
				Token:     packhub.Pointer("token123"),
				Password:  packhub.Pointer("password123"),
			},
			expected: map[string]any{
				"status":     false,
				"profile_id": uint(2),
				"token":      packhub.Pointer("token123"),
				"password":   packhub.Pointer("password123"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.auth.ToMap()
			assert.Equal(t, tt.expected, *result)
		})
	}
}
