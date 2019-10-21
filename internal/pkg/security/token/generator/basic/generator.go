package token

import (
	"context"
	"encoding/json"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	tokenGenerator "github.com/BRBussy/bizzle/internal/pkg/security/token/generator"
	"github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	"gopkg.in/square/go-jose.v2"
)

type generator struct {
	tokenSigner jose.Signer
	validator   validator.Validator
}

func New(
	tokenSigner jose.Signer,
	validator validator.Validator,
) tokenGenerator.Generator {
	return &generator{
		tokenSigner: tokenSigner,
		validator:   validator,
	}
}

func (g *generator) GenerateToken(ctx context.Context, request *tokenGenerator.GenerateTokenRequest) (*tokenGenerator.GenerateTokenResponse, error) {
	if err := g.validator.ValidateRequest(request); err != nil {
		return nil, err
	}

	// marshall claims
	claimsPayload, err := json.Marshal(claims.Serialized{
		Claims: request.Claims,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not marshal claims for token")
		return nil, bizzleException.ErrUnexpected{}
	}

	// sign marshalled payload
	signedObj, err := g.tokenSigner.Sign(claimsPayload)
	if err != nil {
		log.Error().Err(err).Msg("could not sign payload")
		return nil, bizzleException.ErrUnexpected{}
	}

	// serialize signed object
	signedJWT, err := signedObj.CompactSerialize()
	if err != nil {
		log.Error().Err(err).Msg("could not serialize signed token")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &tokenGenerator.GenerateTokenResponse{Token: signedJWT}, nil
}
