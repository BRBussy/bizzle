package jsonRPC

import (
	"context"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
	"net/http"
	"time"
)

type ValidatorJSONRPCClient struct {
	rpcClient jsonrpc.RPCClient
}

func (a *ValidatorJSONRPCClient) Setup(
	urlEndpoint, preSharedSecret string,
) Validator {
	a.rpcClient = jsonrpc.NewClientWithOpts(
		urlEndpoint,
		&jsonrpc.RPCClientOpts{
			HTTPClient:    &http.Client{Timeout: time.Second * 10},
			CustomHeaders: map[string]string{"Pre-Shared-Secret": preSharedSecret},
		},
	)
	return a
}

func (a *ValidatorJSONRPCClient) Validate(ctx context.Context, request ValidateRequest) (*ValidateResponse, error) {
	// perform json rpc request
	rpcResponse, err := a.rpcClient.Call(
		"TokenValidator.Validate",
		&ValidateJSONRPCRequest{
			Token: request.Token,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("TokenValidator.Validate json rpc")
		return nil, exception.ErrUnexpected{}
	}
	if rpcResponse.Error != nil {
		log.Error().Err(rpcResponse.Error).Msg("TokenValidator.Validate json rpc")
		return nil, exception.ErrUnexpected{}
	}

	// parse response
	response := new(ValidateJSONRPCResponse)
	if err := rpcResponse.GetObject(response); err != nil {
		log.Error().Err(err).Msg("parse response object")
		return nil, exception.ErrUnexpected{}
	}

	return &ValidateResponse{MarshalledClaims: response.MarshalledClaims}, nil
}
