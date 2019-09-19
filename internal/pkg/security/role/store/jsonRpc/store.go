package jsonRpc

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	authenticatedJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/authenticated"
	basicJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/basic"
	"github.com/BRBussy/bizzle/internal/pkg/environment"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	roleStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/security/role/store/adaptor/jsonRpc"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

type store struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	env environment.Environment,
	url string,
) (roleStore.Store, error) {
	var client jsonRpcClient.Client
	var err error
	switch env {
	case environment.Development:
		client = basicJsonRpcClient.New(url)
	case environment.Production:
		client, err = authenticatedJsonRpcClient.New(url)
		if err != nil {
			log.Error().Err(err).Msg("creating new authenticated json rpc client")
			return nil, err
		}
	default:
		return nil, bizzleException.ErrUnexpected{Reasons: []string{"invalid environment", env.String()}}
	}
	return &store{
		jsonRpcClient: client,
	}, nil
}

func (a *store) CreateOne(request *roleStore.CreateOneRequest) (*roleStore.CreateOneResponse, error) {
	createResponse := new(roleStoreJsonRpcAdaptor.CreateOneResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		roleStore.CreateOneService,
		roleStoreJsonRpcAdaptor.CreateOneRequest{
			Role: request.Role,
		},
		createResponse); err != nil {
		log.Error().Err(err).Msg("role jsonrpc store create")
		return nil, err
	}

	return &roleStore.CreateOneResponse{}, nil
}

func (a *store) FindOne(request *roleStore.FindOneRequest) (*roleStore.FindOneResponse, error) {
	findOneResponse := new(roleStoreJsonRpcAdaptor.FindOneResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		roleStore.FindOneService,
		roleStoreJsonRpcAdaptor.FindOneRequest{
			Identifier: identifier.Serialized{
				Identifier: request.Identifier,
			},
		},
		findOneResponse); err != nil {
		log.Error().Err(err).Msg("role jsonrpc store find one")
		return nil, err
	}

	return &roleStore.FindOneResponse{
		Role: findOneResponse.Role,
	}, nil
}

func (a *store) UpdateOne(request *roleStore.UpdateOneRequest) (*roleStore.UpdateOneResponse, error) {
	updateOneResponse := new(roleStoreJsonRpcAdaptor.UpdateOneResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		roleStore.UpdateOneService,
		roleStoreJsonRpcAdaptor.UpdateOneRequest{
			Role: request.Role,
		},
		updateOneResponse); err != nil {
		log.Error().Err(err).Msg("role jsonrpc store update one")
		return nil, err
	}

	return &roleStore.UpdateOneResponse{}, nil
}
