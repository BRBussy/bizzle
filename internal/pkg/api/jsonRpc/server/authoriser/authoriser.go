package authoriser

import (
	wrappedClaims "github.com/BRBussy/bizzle/internal/pkg/security/claims/wrapped"
)

type Authoriser interface {
	AuthoriseServiceMethod(jwt string, jsonRpcMethod string) (wrappedClaims.Wrapped, error)
}
