package jsonrpc

import (
	"net/http"

	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntryIgnored "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored"
	budgetEntryIgnoredStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

type adaptor struct {
	store budgetEntryIgnoredStore.Store
}

func New(
	store budgetEntryIgnoredStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetEntryIgnoredStore.ServiceProvider
}

type FindOneRequest struct {
	Identifier identifier.Serialized `json:"identifier"`
}

type FindOneResponse struct {
	Ignored budgetEntryIgnored.Ignored `json:"ignored"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	findOneResponse, err := a.store.FindOne(
		budgetEntryIgnoredStore.FindOneRequest{
			Identifier: request.Identifier.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.Ignored = findOneResponse.Ignored

	return nil
}

type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

type FindManyResponse struct {
	Records []budgetEntryIgnored.Ignored `json:"records"`
	Total   int64                        `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return bizzleException.ErrUnexpected{}
	}

	findOneResponse, err := a.store.FindMany(
		budgetEntryIgnoredStore.FindManyRequest{
			Claims:   c,
			Criteria: request.Criteria.Criteria,
			Query:    request.Query,
		},
	)
	if err != nil {
		return err
	}

	response.Records = findOneResponse.Records
	response.Total = findOneResponse.Total

	return nil
}
