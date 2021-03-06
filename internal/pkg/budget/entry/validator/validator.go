package validator

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
)

// Validator provides validation functions for budget entries
type Validator interface {
	ValidateForCreate(*ValidateForCreateRequest) (*ValidateForCreateResponse, error)
	ValidateForUpdate(*ValidateForUpdateRequest) (*ValidateForUpdateResponse, error)
}

// ValidateForCreateRequest is the request object for the ValidateForCreate service
type ValidateForCreateRequest struct {
	Claims      claims.Claims `validate:"required"`
	BudgetEntry budgetEntry.Entry
}

// ValidateForCreateResponse is the response object for the ValidateForCreate service
type ValidateForCreateResponse struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}

// ValidateForUpdateRequest is the request object for the ValidateForUpdate service
type ValidateForUpdateRequest struct {
	Claims      claims.Claims `validate:"required"`
	BudgetEntry budgetEntry.Entry
}

// ValidateForUpdateResponse is the response object for the ValidateForUpdate service
type ValidateForUpdateResponse struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}
