package validator

import "github.com/go-playground/validator/v10"

type validatorStruct struct {
	validator *validator.Validate
}

var StructValidator *validatorStruct

func init() {
	StructValidator = &validatorStruct{
		validator: validator.New(),
	}
}

func (v validatorStruct) Validate(data interface{}) error {
	if errs := v.validator.Struct(data); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			return &ErrorValidator{err.Field(), err.Tag(), err.Param(), err.Value()}
		}
	}

	return nil
}
