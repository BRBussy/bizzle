package jsonRpc

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/ybbus"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	roleStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/security/role/store/adaptor/jsonRpc"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

type store struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(url, preSharedSecret string) roleStore.Store {
	log.Info().Msg("role json rpc store for: " + url)
	return &store{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (s *store) CreateOne(request *roleStore.CreateOneRequest) (*roleStore.CreateOneResponse, error) {
	createResponse := new(roleStoreJsonRpcAdaptor.CreateOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
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

func (s *store) FindOne(request *roleStore.FindOneRequest) (*roleStore.FindOneResponse, error) {
	findOneResponse := new(roleStoreJsonRpcAdaptor.FindOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
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

func (s *store) FindMany(request *roleStore.FindManyRequest) (*roleStore.FindManyResponse, error) {
	findManyResponse := new(roleStoreJsonRpcAdaptor.FindManyResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		roleStore.FindManyService,
		roleStoreJsonRpcAdaptor.FindManyRequest{
			Criteria: criteria.Serialized{
				Criteria: request.Criteria,
			},
			Query: request.Query,
		},
		findManyResponse); err != nil {
		log.Error().Err(err).Msg("role jsonrpc store find many")
		return nil, err
	}

	return &roleStore.FindManyResponse{
		Records: findManyResponse.Records,
		Total:   findManyResponse.Total,
	}, nil
}

func (s *store) UpdateOne(request *roleStore.UpdateOneRequest) (*roleStore.UpdateOneResponse, error) {
	updateOneResponse := new(roleStoreJsonRpcAdaptor.UpdateOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
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
