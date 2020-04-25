package validator

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
)

type Validator interface {
	ValidateForCreate(ValidateForCreateRequest) (*ValidateForCreateResponse, error)
}

type ValidateForCreateRequest struct {
	Claims claims.Claims `validate:"required"`
	User   user.User     `validate:"required"`
}

type ValidateForCreateResponse struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}
