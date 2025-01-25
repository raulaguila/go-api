package validator

import (
	"fmt"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
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
		input   testStruct
		wantErr bool
	}{
		{"Valid", testStruct{ID: 1, Name: "John", Email: "john@example.com"}, false},
		{"MissingName", testStruct{ID: 1, Name: "", Email: "john@example.com"}, true},
		{"InvalidEmail", testStruct{ID: 1, Name: "John", Email: "john@com"}, true},
		{"NegativeID", testStruct{ID: -1, Name: "John", Email: "john@example.com"}, true},
		{"ZeroID", testStruct{ID: 0, Name: "John", Email: "john@example.com"}, true},
		{"EmptyStruct", testStruct{}, true},
		{"ValidWithLongName", testStruct{ID: 2, Name: "John Doe", Email: "johndoe@example.com"}, false},
		{"VeryLongEmail", testStruct{ID: 3, Name: "Jane", Email: "a.really.long.email.address@example.com"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validatorStruct{validator: validatorInstance}
			err := v.Validate(tt.input)
			assert.Equal(t, tt.wantErr, err != nil, fmt.Sprintf("Test name: %v - Validate() error = %v, wantErr %v", tt.name, err, tt.wantErr))
		})
	}
}
