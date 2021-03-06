package jsonRPC

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/ybbus"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	userStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/user/store/adaptor/jsonRpc"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type store struct {
	jsonRpcClient jsonRpcClient.Client
	validator     validationValidator.Validator
}

func New(
	validator validationValidator.Validator,
	url, preSharedSecret string,
) userStore.Store {
	log.Info().Msg("user json rpc store for: " + url)
	return &store{
		validator:     validator,
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
	}
}

func (s *store) CreateOne(request userStore.CreateOneRequest) (*userStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	createResponse := new(userStoreJsonRpcAdaptor.CreateOneResponse)
	if err := s.jsonRpcClient.JSONRPCRequest(
		userStore.CreateOneService,
		nil,
		userStoreJsonRpcAdaptor.CreateOneRequest{
			User: request.User,
		},
		createResponse,
	); err != nil {
		log.Error().Err(err).Msg("user jsonrpc store create")
		return nil, err
	}
	return &userStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request userStore.FindOneRequest) (*userStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	findOneResponse := new(userStoreJsonRpcAdaptor.FindOneResponse)
	if err := s.jsonRpcClient.JSONRPCRequest(
		userStore.FindOneService,
		request.Claims,
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

func (s *store) FindMany(request userStore.FindManyRequest) (*userStore.FindManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	findManyResponse := new(userStoreJsonRpcAdaptor.FindManyResponse)
	if err := s.jsonRpcClient.JSONRPCRequest(
		userStore.FindManyService,
		request.Claims,
		userStoreJsonRpcAdaptor.FindManyRequest{
			Criteria: criteria.Serialized{
				Criteria: request.Criteria,
			},
			Query: request.Query,
		},
		findManyResponse); err != nil {
		log.Error().Err(err).Msg("user jsonrpc store find many")
		return nil, err
	}

	return &userStore.FindManyResponse{
		Records: findManyResponse.Records,
		Total:   findManyResponse.Total,
	}, nil
}

func (s *store) UpdateOne(request userStore.UpdateOneRequest) (*userStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	updateOneResponse := new(userStoreJsonRpcAdaptor.UpdateOneResponse)
	if err := s.jsonRpcClient.JSONRPCRequest(
		userStore.UpdateOneService,
		request.Claims,
		userStoreJsonRpcAdaptor.UpdateOneRequest{
			User: request.User,
		},
		updateOneResponse); err != nil {
		log.Error().Err(err).Msg("user jsonrpc store update one")
		return nil, err
	}

	return &userStore.UpdateOneResponse{}, nil
}
