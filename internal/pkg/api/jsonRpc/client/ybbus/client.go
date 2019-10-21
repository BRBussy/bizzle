package ybbus

import (
	jsonRPCClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
	"net/http"
	"time"
)

type client struct {
	jsonrpc.RPCClient
}

func New(
	url, preSharedSecret string,
) jsonRPCClient.Client {
	return &client{
		RPCClient: jsonrpc.NewClientWithOpts(
			url,
			&jsonrpc.RPCClientOpts{
				HTTPClient:    &http.Client{Timeout: time.Second * 10},
				CustomHeaders: map[string]string{"Pre-Shared-Secret": preSharedSecret},
			},
		),
	}
}

func (c *client) JsonRpcRequest(method string, request, response interface{}) error {
	// perform json rpc request
	rpcResponse, err := c.Call(
		method,
		[]interface{}{request},
	)
	if err != nil {
		return err
	}
	if rpcResponse.Error != nil {
		return rpcResponse.Error
	}

	// parse response
	if err := rpcResponse.GetObject(response); err != nil {
		log.Error().Err(err).Msg("parse response object")
		return err
	}

	return nil
}
