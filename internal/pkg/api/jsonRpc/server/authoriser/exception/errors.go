package exception

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/security/permission"
)

type NotAuthorised struct {
	Permission permission.Permission
}

func (e NotAuthorised) Error() string {
	return fmt.Sprintf("not authorised for %s", e.Permission)
}

type InvalidClaims struct {
	ExpectedClaimsType claims.Type
}

func (e InvalidClaims) Error() string {
	return fmt.Sprintf("invalid claims, expected: %s", e.ExpectedClaimsType)
}
