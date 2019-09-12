package jsonRpc

import (
	jsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client"
	authenticatedJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/authenticated"
	basicJsonRpcClient "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/client/basic"
	"github.com/BRBussy/bizzle/internal/pkg/environment"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	exerciseStoreJsonRpcAdaptor "github.com/BRBussy/bizzle/internal/pkg/exercise/store/adaptor/jsonRpc"
	wrappedCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/wrapped"
	"github.com/rs/zerolog/log"
)

type store struct {
	jsonRpcClient jsonRpcClient.Client
}

func New(
	env environment.Environment,
	authenticatorURL string,
) (exerciseStore.Store, error) {
	var client jsonRpcClient.Client
	var err error
	switch env {
	case environment.Development:
		client = basicJsonRpcClient.New(authenticatorURL)
	case environment.Production:
		client, err = authenticatedJsonRpcClient.New(authenticatorURL)
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

func (a *store) Find(request *exerciseStore.FindRequest) (*exerciseStore.FindResponse, error) {
	wrappedCriteria := make([]wrappedCriterion.Wrapped, 0)
	for _, crit := range request.Criteria {
		wrappedCrit, err := wrappedCriterion.Wrap(crit)
		if err != nil {
			log.Error().Err(err).Msg("wrapping criterion")
			return nil, err
		}
		wrappedCriteria = append(wrappedCriteria, *wrappedCrit)
	}

	findResponse := new(exerciseStoreJsonRpcAdaptor.FindResponse)
	if err := a.jsonRpcClient.JsonRpcRequest(
		exerciseStore.FindService,
		exerciseStoreJsonRpcAdaptor.FindRequest{},
		findResponse); err != nil {
		log.Error().Err(err).Msg("authenticator json rpc SignUp")
		return nil, err
	}

	return &exerciseStore.FindResponse{}, nil
}
