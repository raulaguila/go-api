package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// validatorStruct is a type that encapsulates a validator.Validate instance for struct validation operations.
type validatorStruct struct {
	validator *validator.Validate
}

// StructValidator is a global instance of validatorStruct used to perform data structure validation.
var StructValidator *validatorStruct

// init initializes the StructValidator with a new instance of validatorStruct containing a new validator.
func init() {
	StructValidator = &validatorStruct{
		validator: validator.New(),
	}
}

// Validate checks if the input data satisfies all validation rules defined in the schema.
// Returns a ValidateError containing details of the validation failure, or nil if the data is valid.
// Utilizes a validator.Validate instance from the validatorStruct to perform struct-level data validation.
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
