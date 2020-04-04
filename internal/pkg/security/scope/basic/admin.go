package basic

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/scope"
	"github.com/BRBussy/bizzle/pkg/search/criterion/text"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

// Admin is the scope.Admin - used to determine scope
type Admin struct {
	validator validationValidator.Validator
}

// New creates a new scope.Manager
func New(
	validator validationValidator.Validator,
) scope.Admin {
	return &Admin{
		validator: validator,
	}
}

// ApplyScopeToCriteria is used to add scope from claims to a search criteria
func (a *Admin) ApplyScopeToCriteria(request scope.ApplyScopeToCriteriaRequest) (*scope.ApplyScopeToCriteriaResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &scope.ApplyScopeToCriteriaResponse{
		ScopedCriteria: append(
			request.CriteriaToScope,
			text.Exact{
				Field: "ownerID",
				Text:  request.Claims.ScopingID().String(),
			}),
	}, nil
}

// ApplyScopeToIdentifier is used to add scope from claims to a search identifier
func (a *Admin) ApplyScopeToIdentifier(request scope.ApplyScopeToIdentifierRequest) (*scope.ApplyScopeToIdentifierResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &scope.ApplyScopeToIdentifierResponse{
		ScopedIdentifier: scope.ScopedIdentifier{
			IdentifierToScope: request.IdentifierToScope,
			OwnerID:           request.Claims.ScopingID(),
		},
	}, nil
}
