package jsonRPC

import (
	"net/http"

	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	budgetCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/rs/zerolog/log"
)

type adaptor struct {
	store budgetCategoryRuleStore.Store
}

// New creates a new budget category rule store json rpc adaptor
func New(
	store budgetCategoryRuleStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetCategoryRuleStore.ServiceProvider
}

// FindManyRequest is the request object for FindMany method
type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

// FindManyResponse is the response object for the FindMany method
type FindManyResponse struct {
	Records []budgetEntryCategoryRule.CategoryRule `json:"records"`
	Total   int64                                  `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return bizzleException.ErrUnexpected{}
	}

	findManyResponse, err := a.store.FindMany(&budgetCategoryRuleStore.FindManyRequest{
		Claims:   c,
		Criteria: request.Criteria.Criteria,
		Query:    request.Query,
	})
	if err != nil {
		return err
	}

	response.Records = findManyResponse.Records
	response.Total = findManyResponse.Total

	return nil
}
