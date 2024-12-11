package validator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateError_Error(t *testing.T) {
	tests := []struct {
		name      string
		err       ValidateError
		wantError string
	}{
		{
			name: "with param",
			err: ValidateError{
				field: "age",
				tag:   "min",
				param: "18",
				value: 16,
			},
			wantError: `{"value":"16","field":"age","tag":"min","param":"18"}`,
		},
		{
			name: "without param",
			err: ValidateError{
				field: "email",
				tag:   "required",
				value: "",
			},
			wantError: `{"value":"","field":"email","tag":"required","param":""}`,
		},
		{
			name: "complex value",
			err: ValidateError{
				field: "details",
				tag:   "type",
				param: "map",
				value: map[string]interface{}{"key": "value"},
			},
			wantError: `{"value":"map[key:value]","field":"details","tag":"type","param":"map"}`,
		},
		{
			name:      "empty values",
			err:       ValidateError{},
			wantError: `{"value":"<nil>","field":"","tag":"","param":""}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			assert.Equal(t, tt.wantError, got)
		})
	}
}
