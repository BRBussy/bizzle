package jsonRpc

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	brizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	brizzleAuthenticatorJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	"github.com/rs/zerolog/log"
)

type authenticator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) brizzleAuthenticator.Authenticator {
	return &authenticator{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *authenticator) SignUp(*brizzleAuthenticator.SignUpRequest) (*brizzleAuthenticator.SignUpResponse, error) {
	signUpResponse := new(brizzleAuthenticatorJsonRpcAdaptor.SignUpResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		brizzleAuthenticator.SignUpService,
		brizzleAuthenticatorJsonRpcAdaptor.SignUpRequest{},
		signUpResponse); err != nil {
		log.Error().Err(err).Msg("authenticator json rpc SignUp")
		return nil, err
	}

	return &brizzleAuthenticator.SignUpResponse{
		Msg: signUpResponse.Msg,
	}, nil
}
