package validator

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

type validatorStruct struct {
	validator *validator.Validate
}

var StructValidator *validatorStruct

func init() {
	StructValidator = &validatorStruct{
		validator: validator.New(),
	}
}

func (v validatorStruct) Validate(data any) error {
	if result := v.validator.Struct(data); result != nil {
		var errs validator.ValidationErrors
		if errors.As(errs, &errs) {
			for _, err := range errs {
				return &ValidateError{err.Field(), err.Tag(), err.Param(), err.Value()}
			}
		}
	}

	return nil
}
