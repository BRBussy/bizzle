package jsonRPC

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/ybbus"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	userStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/user/store/adaptor/jsonRpc"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

type store struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	url, preSharedSecret string,
) userStore.Store {
	log.Info().Msg("user json rpc store for: " + url)
	return &store{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (s *store) CreateOne(request *userStore.CreateOneRequest) (*userStore.CreateOneResponse, error) {
	createResponse := new(userStoreJsonRpcAdaptor.CreateOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		userStore.CreateOneService,
		userStoreJsonRpcAdaptor.CreateOneRequest{
			User: request.User,
		},
		createResponse); err != nil {
		log.Error().Err(err).Msg("user jsonrpc store create")
		return nil, err
	}
	return &userStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request *userStore.FindOneRequest) (*userStore.FindOneResponse, error) {
	findOneResponse := new(userStoreJsonRpcAdaptor.FindOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		userStore.FindOneService,
		userStoreJsonRpcAdaptor.FindOneRequest{
			Identifier: identifier.Serialized{
				Identifier: request.Identifier,
			},
		},
		findOneResponse); err != nil {
		log.Error().Err(err).Msg("user jsonrpc store find one")
		return nil, err
	}

	return &userStore.FindOneResponse{
		User: findOneResponse.User,
	}, nil
}

func (s *store) UpdateOne(request *userStore.UpdateOneRequest) (*userStore.UpdateOneResponse, error) {
	updateOneResponse := new(userStoreJsonRpcAdaptor.UpdateOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		userStore.UpdateOneService,
		userStoreJsonRpcAdaptor.UpdateOneRequest{
			User: request.User,
		},
		updateOneResponse); err != nil {
		log.Error().Err(err).Msg("user jsonrpc store update one")
		return nil, err
	}

	return &userStore.UpdateOneResponse{}, nil
}
