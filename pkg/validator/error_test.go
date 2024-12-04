package validator

import (
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
			wantError: "value: '16' in the 'age' field does not meet the requirement: min[18]",
		},
		{
			name: "without param",
			err: ValidateError{
				field: "email",
				tag:   "required",
				value: "",
			},
			wantError: "value: '' in the 'email' field does not meet the requirement: required",
		},
		{
			name: "complex value",
			err: ValidateError{
				field: "details",
				tag:   "type",
				param: "map",
				value: map[string]interface{}{"key": "value"},
			},
			wantError: "value: 'map[key:value]' in the 'details' field does not meet the requirement: type[map]",
		},
		{
			name:      "empty values",
			err:       ValidateError{},
			wantError: "value: '<nil>' in the '' field does not meet the requirement: ",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.wantError {
				t.Errorf("got %v, want %v", got, tt.wantError)
			}
		})
	}
}
