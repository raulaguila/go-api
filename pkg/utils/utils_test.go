package utils

import (
	"errors"
	"testing"
)

func TestPanicIfErr(t *testing.T) {
	tests := []struct {
		name        string
		err         error
		shouldPanic bool
	}{
		{"nilError", nil, false},
		{"nonNilError", errors.New("some error"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil && !tt.shouldPanic {
					t.Errorf("PanicIfErr() panicked when it shouldn't: %v", r)
				}
			}()
			PanicIfErr(tt.err)
		})
	}
}
