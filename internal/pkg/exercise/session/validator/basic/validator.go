package basic

import (
	sessionValidator "github.com/BRBussy/bizzle/internal/pkg/exercise/session/validator"
)

type validator struct {
}

func New() sessionValidator.Validator {
	return &validator{}
}

func (v *validator) ValidateForCreate(request *sessionValidator.ValidateForCreateRequest) (*sessionValidator.ValidateForCreateResponse, error) {
	return &sessionValidator.ValidateForCreateResponse{}, nil
}
