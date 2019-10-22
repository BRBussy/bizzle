package jsonRPC

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/ybbus"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	tokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator"
	tokenValidatorJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/security/token/validator/adaptor/jsonRPC"
	"github.com/rs/zerolog/log"
)

type validator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	url, preSharedSecret string,
) tokenValidator.Validator {
	return &validator{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (a *validator) Validate(request *tokenValidator.ValidateRequest) (*tokenValidator.ValidateResponse, error) {
	response := new(tokenValidatorJSONRPCAdaptor.ValidateResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		tokenValidator.ValidateService,
		tokenValidatorJSONRPCAdaptor.ValidateRequest{
			Token: request.Token,
		},
		response,
	); err != nil {
		log.Error().Err(err).Msg("TokenValidator.Validate json rpc")
		return nil, exception.ErrUnexpected{}
	}
	return &tokenValidator.ValidateResponse{MarshalledClaims: response.MarshalledClaims}, nil
}
