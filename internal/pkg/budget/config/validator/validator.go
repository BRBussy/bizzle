package validator

import (
	budgetConfig "github.com/BRBussy/bizzle/internal/pkg/budget/config"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
)

// Validator provides validation functions for budget config
type Validator interface {
	ValidateForCreate(*ValidateForCreateRequest) (*ValidateForCreateResponse, error)
	ValidateForUpdate(*ValidateForUpdateRequest) (*ValidateForUpdateResponse, error)
}

// ValidateForCreateRequest is the request object for the ValidateForCreate service
type ValidateForCreateRequest struct {
	Claims       claims.Claims `validate:"required"`
	BudgetConfig budgetConfig.Config
}

// ValidateForCreateResponse is the response object for the ValidateForCreate service
type ValidateForCreateResponse struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}

// ValidateForUpdateRequest is the request object for the ValidateForUpdate service
type ValidateForUpdateRequest struct {
	Claims       claims.Claims `validate:"required"`
	BudgetConfig budgetConfig.Config
}

// ValidateForUpdateResponse is the response object for the ValidateForUpdate service
type ValidateForUpdateResponse struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}
