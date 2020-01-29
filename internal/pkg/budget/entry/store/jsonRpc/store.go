package jsonRpc

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	ybbusJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/ybbus"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	budgetEntryStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store/adaptor/jsonRpc"
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
) budgetEntryStore.Store {
	log.Info().Msg("budgetEntry json rpc store for: " + url)
	return &store{
		jsonRpcClient: ybbusJsonRpcClient.New(url, preSharedSecret),
		validator:     validator,
	}
}

func (s *store) CreateOne(request *budgetEntryStore.CreateOneRequest) (*budgetEntryStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	createResponse := new(budgetEntryStoreJsonRpcAdaptor.CreateOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		budgetEntryStore.CreateOneService,
		budgetEntryStoreJsonRpcAdaptor.CreateOneRequest{
			Entry: request.Entry,
		},
		createResponse); err != nil {
		log.Error().Err(err).Msg("budgetEntry jsonrpc store create")
		return nil, err
	}
	return &budgetEntryStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request *budgetEntryStore.FindOneRequest) (*budgetEntryStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	findOneResponse := new(budgetEntryStoreJsonRpcAdaptor.FindOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		budgetEntryStore.FindOneService,
		budgetEntryStoreJsonRpcAdaptor.FindOneRequest{
			Identifier: identifier.Serialized{
				Identifier: request.Identifier,
			},
		},
		findOneResponse); err != nil {
		log.Error().Err(err).Msg("budgetEntry jsonrpc store find one")
		return nil, err
	}

	return &budgetEntryStore.FindOneResponse{
		Entry: findOneResponse.Entry,
	}, nil
}

func (s *store) FindMany(request *budgetEntryStore.FindManyRequest) (*budgetEntryStore.FindManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	findManyResponse := new(budgetEntryStoreJsonRpcAdaptor.FindManyResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		budgetEntryStore.FindManyService,
		budgetEntryStoreJsonRpcAdaptor.FindManyRequest{
			Criteria: criteria.Serialized{
				Criteria: request.Criteria,
			},
			Query: request.Query,
		},
		findManyResponse); err != nil {
		log.Error().Err(err).Msg("budgetEntry jsonrpc store find many")
		return nil, err
	}

	return &budgetEntryStore.FindManyResponse{
		Records: findManyResponse.Records,
		Total:   findManyResponse.Total,
	}, nil
}

func (s *store) UpdateOne(request *budgetEntryStore.UpdateOneRequest) (*budgetEntryStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	updateOneResponse := new(budgetEntryStoreJsonRpcAdaptor.UpdateOneResponse)
	if err := s.jsonRpcClient.JsonRpcRequest(
		budgetEntryStore.UpdateOneService,
		budgetEntryStoreJsonRpcAdaptor.UpdateOneRequest{
			Entry: request.Entry,
		},
		updateOneResponse); err != nil {
		log.Error().Err(err).Msg("budgetEntry jsonrpc store update one")
		return nil, err
	}

	return &budgetEntryStore.UpdateOneResponse{}, nil
}
