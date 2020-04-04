package jsonrpc

import (
	"net/http"

	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetConfig "github.com/BRBussy/bizzle/internal/pkg/budget/config"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/config/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

type adaptor struct {
	store budgetEntryStore.Store
}

func New(
	store budgetEntryStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetEntryStore.ServiceProvider
}

type CreateOneRequest struct {
	Config budgetConfig.Config `json:"budgetConfig"`
}

type CreateOneResponse struct {
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		budgetEntryStore.CreateOneRequest{
			Config: request.Config,
		},
	); err != nil {
		return err
	}

	return nil
}

type FindOneRequest struct {
	Identifier identifier.Serialized `json:"identifier"`
}

type FindOneResponse struct {
	Config budgetConfig.Config `json:"budgetConfig"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	findOneResponse, err := a.store.FindOne(
		budgetEntryStore.FindOneRequest{
			Identifier: request.Identifier.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.Config = findOneResponse.Config

	return nil
}

type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

type FindManyResponse struct {
	Records []budgetConfig.Config `json:"records"`
	Total   int64                 `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return bizzleException.ErrUnexpected{}
	}

	findOneResponse, err := a.store.FindMany(
		budgetEntryStore.FindManyRequest{
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

type UpdateOneRequest struct {
	Config budgetConfig.Config `json:"budgetConfig"`
}

type UpdateOneResponse struct {
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	if _, err := a.store.UpdateOne(
		budgetEntryStore.UpdateOneRequest{
			Config: request.Config,
		},
	); err != nil {
		return err
	}
	return nil
}
