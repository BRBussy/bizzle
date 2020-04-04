package scope

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

// Admin is a service provider that provides scoping using claims
type Admin interface {
	ApplyScopeToCriteria(ApplyScopeToCriteriaRequest) (*ApplyScopeToCriteriaResponse, error)
	ApplyScopeToIdentifier(ApplyScopeToIdentifierRequest) (*ApplyScopeToIdentifierResponse, error)
}

// ApplyScopeToCriteriaRequest is the request object for the ApplyScopeToCriteria Service
type ApplyScopeToCriteriaRequest struct {
	Claims          claims.Claims     `validate:"required"`
	CriteriaToScope criteria.Criteria `validate:"required"`
}

// ApplyScopeToCriteriaResponse is the response object for the ApplyScopeToCriteria Service
type ApplyScopeToCriteriaResponse struct {
	ScopedCriteria criteria.Criteria `validate:"required"`
}

// ApplyScopeToIdentifierRequest is the request object for the ApplyScopeToIdentifier Service
type ApplyScopeToIdentifierRequest struct {
	Claims            claims.Claims         `validate:"required"`
	IdentifierToScope identifier.Identifier `validate:"required"`
}

// ApplyScopeToIdentifierResponse is the response object for the ApplyScopeToIdentifier Service
type ApplyScopeToIdentifierResponse struct {
	ScopedIdentifier identifier.Identifier
}
