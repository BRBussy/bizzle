package request

import (
	"fmt"
	validateValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	goValidator "gopkg.in/go-playground/validator.v9"
)

type validator struct {
	validate *goValidator.Validate
}

func New() validateValidator.Validator {
	return &validator{
		validate: goValidator.New(),
	}
}

func (v validator) Validate(request interface{}) error {
	var reasons []string

	if err := v.validate.Struct(request); err != nil {
		validationErrors := err.(goValidator.ValidationErrors)
		for key, value := range validationErrors {
			reasons = append(reasons, fmt.Sprintf("'%d' %s", key, value))
		}
	}

	if len(reasons) > 0 {
		return validateValidator.ErrRequestNotValid{Reasons: reasons}
	}
	return nil
}
