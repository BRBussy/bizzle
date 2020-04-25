package client

import "github.com/BRBussy/bizzle/internal/pkg/security/claims"

type Client interface {
	JSONRPCRequest(method string, authClaims claims.Claims, request, response interface{}) error
}
