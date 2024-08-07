package validator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type structTest struct {
	Name  string `validate:"required,min=5,max=10"`
	Age   int    `validate:"required,min=12,max=18"`
	Email string `validate:"required,email"`
}

// go test -run TestValidatorWithoutData
func TestValidatorWithoutData(t *testing.T) {
	element := &structTest{}

	err := StructValidator.Validate(element)
	require.Error(t, err)
	require.ErrorAs(t, err, &ErrValidator)
}

// go test -run TestValidatorWitInvalidData
func TestValidatorWitInvalidData(t *testing.T) {
	element := &structTest{"1234", 22, "toError@email"}

	err := StructValidator.Validate(element)
	require.Error(t, err)
	require.ErrorAs(t, err, &ErrValidator)

	element.Name = "0123456"
	element.Age = 33
	element.Email = "email.com"

	err = StructValidator.Validate(element)
	require.Error(t, err)
	require.IsType(t, "", err.Error())
	require.ErrorAs(t, err, &ErrValidator)
}

// go test -run TestValidatorWitValidData
func TestValidatorWitValidData(t *testing.T) {
	element := &structTest{"123456", 15, "example@example.com"}

	err := StructValidator.Validate(element)
	require.NoError(t, err)

	element.Name = "12345678"
	element.Age = 13
	element.Email = "email@email.com"

	err = StructValidator.Validate(element)
	require.NoError(t, err)
}
