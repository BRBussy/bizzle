package jsonRpc

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	brizzleAuthenticatorJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/authenticator/adaptor/jsonRpc"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	"github.com/rs/zerolog/log"
)

type store struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	jsonRpcClient jsonRpcClient.Client,
) exerciseStore.Store {
	return &store{
		jsonRpcClient: jsonRpcClient,
	}
}

func (a *store) Find(request *exerciseStore.FindRequest) (*exerciseStore.FindResponse, error) {
	signUpResponse := new(brizzleAuthenticatorJsonRpcAdaptor.SignUpResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		exerciseStore.FindService,
		brizzleAuthenticatorJsonRpcAdaptor.SignUpRequest{},
		signUpResponse); err != nil {
		log.Error().Err(err).Msg("authenticator json rpc SignUp")
		return nil, err
	}

	return &exerciseStore.FindResponse{}, nil
}
