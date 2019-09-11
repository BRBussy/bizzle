package jsonRpc

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	authenticatedJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/authenticated"
	basicJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/basic"
	brizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	brizzleAuthenticatorJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	"github.com/BRBussy/bizzle/internal/pkg/environment"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/rs/zerolog/log"
)

type authenticator struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	env environment.Environment,
	authenticatorURL string,
) (brizzleAuthenticator.Authenticator, error) {
	var client jsonRpcClient.Client
	var err error
	switch env {
	case environment.Production:
		client = basicJsonRpcClient.New(authenticatorURL)
	case environment.Development:
		client, err = authenticatedJsonRpcClient.New(authenticatorURL)
		if err != nil {
			log.Error().Err(err).Msg("creating new authenticated json rpc client")
			return nil, err
		}
	default:
		return nil, bizzleException.ErrUnexpected{Reasons: []string{"invalid environment", env.String()}}
	}

	return &authenticator{
		jsonRpcClient: client,
	}, nil
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
