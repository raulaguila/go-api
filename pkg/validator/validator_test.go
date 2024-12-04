package validator

import (
	"github.com/go-playground/validator/v10"
	"testing"
)

type testStruct struct {
	ID    int    `validate:"gt=0"`
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
}

func TestValidatorStruct_Validate(t *testing.T) {
	validatorInstance := validator.New()
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{"Valid", testStruct{ID: 1, Name: "John", Email: "john@example.com"}, false},
		{"MissingName", testStruct{ID: 1, Name: "", Email: "john@example.com"}, true},
		{"InvalidEmail", testStruct{ID: 1, Name: "John", Email: "john@com"}, true},
		{"NegativeID", testStruct{ID: -1, Name: "John", Email: "john@example.com"}, true},
		{"ZeroID", testStruct{ID: 0, Name: "John", Email: "john@example.com"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validatorStruct{validator: validatorInstance}
			if err := v.Validate(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
