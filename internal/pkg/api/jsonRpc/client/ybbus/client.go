package ybbus

import (
	"encoding/json"
	jsonRPCClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/rs/zerolog/log"
	"github.com/ybbus/jsonrpc"
	"net/http"
	"time"
)

type client struct {
	preSharedSecret string
	url             string
}

func New(
	url, preSharedSecret string,
) jsonRPCClient.Client {
	return &client{
		preSharedSecret: preSharedSecret,
		url:             url,
	}
}

func (c *client) JSONRPCRequest(method string, authClaims claims.Claims, request, response interface{}) error {
	var marshalledClaimsForHeader string
	if authClaims != nil {
		marshalledClaims, err := json.Marshal(claims.Serialized{Claims: authClaims})
		if err != nil {
			log.Error().Err(err).Msg("could not marshall claims")
			return bizzleException.ErrUnexpected{}
		}
		marshalledClaimsForHeader = string(marshalledClaims)
	}

	rpcResponse, err := jsonrpc.NewClientWithOpts(
		c.url,
		&jsonrpc.RPCClientOpts{
			HTTPClient: &http.Client{Timeout: time.Second * 10},
			CustomHeaders: map[string]string{
				"Pre-Shared-Secret": c.preSharedSecret,
				"Claims":            marshalledClaimsForHeader,
			},
		},
	).Call(method, request)

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
