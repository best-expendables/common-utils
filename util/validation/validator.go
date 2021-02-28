package validation

import validCheck "gopkg.in/go-playground/validator.v9"

type Validator interface {
	Validate(s interface{}) error
}

func NewInternalValidator(validator *validCheck.Validate) Validator {
	return validatorImpl{validator: validator}
}

type validatorImpl struct {
	validator *validCheck.Validate
}

func (v validatorImpl) Validate(s interface{}) error {
	return v.validator.Struct(s)
}
