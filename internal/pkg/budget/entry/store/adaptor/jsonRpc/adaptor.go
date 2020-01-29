package jsonRpc

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"net/http"
)

type adaptor struct {
	store budgetEntryStore.Store
}

func New(
	store budgetEntryStore.Store,
) jsonRpcServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return budgetEntryStore.ServiceProvider
}

type CreateOneRequest struct {
	Entry budgetEntry.Entry `json:"budgetEntry"`
}

type CreateOneResponse struct {
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		&budgetEntryStore.CreateOneRequest{
			Entry: request.Entry,
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
	Entry budgetEntry.Entry `json:"budgetEntry"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	findOneResponse, err := a.store.FindOne(
		&budgetEntryStore.FindOneRequest{
			Identifier: request.Identifier.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.Entry = findOneResponse.Entry

	return nil
}

type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

type FindManyResponse struct {
	Records []budgetEntry.Entry `json:"records"`
	Total   int64               `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	findOneResponse, err := a.store.FindMany(
		&budgetEntryStore.FindManyRequest{
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
	Entry budgetEntry.Entry `json:"budgetEntry"`
}

type UpdateOneResponse struct {
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	if _, err := a.store.UpdateOne(
		&budgetEntryStore.UpdateOneRequest{
			Entry: request.Entry,
		},
	); err != nil {
		return err
	}
	return nil
}
