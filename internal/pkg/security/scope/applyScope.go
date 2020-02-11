package scope

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/criterion/text"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

// ApplyScopeToCriteria is used to add scope from claims to a search criteria
func ApplyScopeToCriteria(claimsToApply claims.Claims, criteriaToApplyTo criteria.Criteria) (criteria.Criteria, error) {
	return append(
		criteriaToApplyTo,
		text.Exact{
			Field: "ownerID",
			Text: claimsToApply.ScopingID().String()
		}),
		nil
}

// ApplyScopeToIdentifier is used to add scope from claims to a search identifier
func ApplyScopeToIdentifier(claimsToApply claims.Claims, identifierToApplyTo identifier.Identifier) (identifier.Identifier, error) {

}
