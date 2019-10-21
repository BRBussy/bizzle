package token

import (
	"context"
	"crypto/rsa"
	"github.com/BRBussy/bizzle/internal/pkg/security/token"
	tokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator"
	validateValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"gopkg.in/square/go-jose.v2"
)

type validator struct {
	rsaKeyPair       *rsa.PrivateKey
	requestValidator validateValidator.Validator
}

func New(
	rsaKeyPair *rsa.PrivateKey,
	requestValidator validateValidator.Validator,
) tokenValidator.Validator {
	return &validator{
		rsaKeyPair:       rsaKeyPair,
		requestValidator: requestValidator,
	}
}

func (v *validator) Validate(ctx context.Context, request tokenValidator.ValidateRequest) (*tokenValidator.ValidateResponse, error) {
	if err := v.requestValidator.ValidateRequest(request); err != nil {
		return nil, err
	}

	// Parse the jwt. Successful parse means the received token string is a jwt
	jwtObject, err := jose.ParseSigned(request.Token)
	if err != nil {
		return nil, token.ErrInvalidToken{Reasons: []string{err.Error()}}
	}

	// Verify jwt signature and retrieve json marshalled claims
	// Failure indicates jwt was damaged or tampered with
	jsonClaims, err := jwtObject.Verify(&v.rsaKeyPair.PublicKey)
	if err != nil {
		return nil, token.ErrTokenVerification{Reasons: []string{err.Error()}}
	}

	// return marshalled claims
	return &tokenValidator.ValidateResponse{MarshalledClaims: jsonClaims}, nil
}
