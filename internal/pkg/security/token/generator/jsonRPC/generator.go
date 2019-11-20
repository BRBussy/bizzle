package jsonRPC

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/ybbus"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	tokenGenerator "github.com/BRBussy/bizzle/internal/pkg/security/token/generator"
	tokenGeneratorJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/security/token/generator/adaptor/jsonRPC"
	"github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type generator struct {
	jsonRpcClient jsonRpcClient.Client
	validator     validator.Validator
}

func New(
	url, preSharedSecret string,
	validator validator.Validator,
) tokenGenerator.Generator {
	return &generator{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
		validator:     validator,
	}
}

func (a *generator) GenerateToken(request *tokenGenerator.GenerateTokenRequest) (*tokenGenerator.GenerateTokenResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	generateResponse := new(tokenGeneratorJSONRPCAdaptor.GenerateTokenResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		tokenGenerator.GenerateTokenService,
		tokenGeneratorJSONRPCAdaptor.GenerateTokenRequest{
			Claims: claims.Serialized{
				Claims: request.Claims,
			},
		},
		generateResponse,
	); err != nil {
		log.Error().Err(err).Msg("token jsonrpc generator generate")
		return nil, err
	}
	return &tokenGenerator.GenerateTokenResponse{Token: generateResponse.Token}, nil
}
