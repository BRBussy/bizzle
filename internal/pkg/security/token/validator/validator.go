package validator

import "github.com/BRBussy/bizzle/internal/pkg/security/claims"

type Validator interface {
	Validate(request ValidateRequest) (*ValidateResponse, error)
}

const ServiceProvider = "Token-Validator"

const ValidateService = ServiceProvider + ".Validate"

type ValidateRequest struct {
	Token string `validate:"required"`
}

type ValidateResponse struct {
	Claims claims.Claims
}
