package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// validatorStruct is a structure that holds a validator instance used for validating data structures.
type validatorStruct struct {
	validator *validator.Validate
}

// StructValidator is a global variable that holds a pointer to a validatorStruct instance.
// It is used for validating data structures according to specified validation rules.
var StructValidator *validatorStruct

// init initializes the StructValidator with a new instance of validatorStruct containing a new validator instance.
func init() {
	StructValidator = &validatorStruct{
		validator: validator.New(),
	}
}

// Validate performs a validation check on the provided data using the validator instance within the struct.
// If the validation fails due to constraints violations, it returns a ValidateError with details about the failure.
// The method returns nil if no validation errors are found.
func (v validatorStruct) Validate(data any) error {
	if result := v.validator.Struct(data); result != nil {
		var errs validator.ValidationErrors
		if errors.As(result, &errs) {
			for _, err := range errs {
				return &ValidateError{err.Field(), err.Tag(), err.Param(), err.Value()}
			}
		}
	}

	return nil
}
