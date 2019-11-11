package jsonRpc

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"net/http"
)

type adaptor struct {
	store exerciseStore.Store
}

func New(
	store exerciseStore.Store,
) jsonRpcServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return exerciseStore.ServiceProvider
}

type CreateOneRequest struct {
	Exercise exercise.Serialized `json:"exercise"`
}

type CreateOneResponse struct {
	Exercise exercise.Serialized `json:"exercise"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		&exerciseStore.CreateOneRequest{
			Exercise: request.Exercise.Exercise,
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
	Exercise exercise.Exercise `json:"exercise"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	findOneResponse, err := a.store.FindOne(
		&exerciseStore.FindOneRequest{
			Identifier: request.Identifier.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.Exercise = findOneResponse.Exercise

	return nil
}

type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

type FindManyResponse struct {
	Records []exercise.Exercise `json:"records"`
	Total   int64               `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	findOneResponse, err := a.store.FindMany(
		&exerciseStore.FindManyRequest{
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
	Exercise exercise.Exercise `json:"exercise"`
}

type UpdateOneResponse struct {
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	if _, err := a.store.UpdateOne(
		&exerciseStore.UpdateOneRequest{
			Exercise: request.Exercise,
		},
	); err != nil {
		return err
	}
	return nil
}
