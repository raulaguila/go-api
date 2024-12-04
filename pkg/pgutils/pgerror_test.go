package pgutils

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
)

func TestHandlerError(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected error
	}{
		{
			name:     "Nil Error",
			input:    nil,
			expected: nil,
		},
		{
			name:     "Duplicated Key Error",
			input:    &pgconn.PgError{Code: "23505"},
			expected: ErrDuplicatedKey,
		},
		{
			name:     "Foreign Key Violated Error",
			input:    &pgconn.PgError{Code: "23503"},
			expected: ErrForeignKeyViolated,
		},
		{
			name:     "Undefined Column Error",
			input:    &pgconn.PgError{Code: "42703"},
			expected: ErrUndefinedColumn,
		},
		{
			name:     "Database Already Exists Error",
			input:    &pgconn.PgError{Code: "42P04"},
			expected: ErrDatabaseAlreadyExists,
		},
		{
			name:     "Unknown PgError",
			input:    &pgconn.PgError{Code: "99999"},
			expected: &pgconn.PgError{Code: "99999"},
		},
		{
			name:     "Non PgError",
			input:    errors.New("non pg error"),
			expected: errors.New("non pg error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := HandlerError(tt.input)
			if result != nil && tt.expected != nil && result.Error() != tt.expected.Error() {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}
