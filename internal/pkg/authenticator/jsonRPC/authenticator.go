package jsonRPC

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/ybbus"
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	bizzleAuthenticatorJSONRPCAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type authenticator struct {
	jsonRpcClient jsonRpcClient.Client
	validator     validationValidator.Validator
}

func New(
	validator validationValidator.Validator,
	url, preSharedSecret string,
) bizzleAuthenticator.Authenticator {
	log.Info().Msg("role json rpc store for: " + url)
	return &authenticator{
		validator:     validator,
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (a *authenticator) Login(request *bizzleAuthenticator.LoginRequest) (*bizzleAuthenticator.LoginResponse, error) {
	return nil, bizzleException.ErrUnexpected{Reasons: []string{"not implemented"}}
}

func (a *authenticator) AuthenticateService(request *bizzleAuthenticator.AuthenticateServiceRequest) (*bizzleAuthenticator.AuthenticateServiceResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	authenticateServiceResponse := new(bizzleAuthenticatorJSONRPCAdaptor.AuthenticateServiceResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		bizzleAuthenticator.AuthenticateServiceService,
		bizzleAuthenticatorJSONRPCAdaptor.AuthenticateServiceRequest{
			Claims: claims.Serialized{
				Claims: request.Claims,
			},
			Service: request.Service,
		},
		authenticateServiceResponse); err != nil {
		log.Error().Err(err).Msg("auth authenticator jsonrpc authenticateService")
		return nil, err
	}

	return &bizzleAuthenticator.AuthenticateServiceResponse{}, nil
}
