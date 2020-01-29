package jsonRpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/exercise/session"
	sessionStore "github.com/BRBussy/bizzle/internal/pkg/exercise/session/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"net/http"
)

type adaptor struct {
	store sessionStore.Store
}

func New(
	store sessionStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return sessionStore.ServiceProvider
}

type CreateOneRequest struct {
	Session session.Session `json:"session"`
}

type CreateOneResponse struct {
	Session session.Session `json:"session"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		&sessionStore.CreateOneRequest{
			Session: request.Session,
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
	Session session.Session `json:"session"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	findOneResponse, err := a.store.FindOne(
		&sessionStore.FindOneRequest{
			Identifier: request.Identifier.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.Session = findOneResponse.Session

	return nil
}

type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

type FindManyResponse struct {
	Records []session.Session `json:"records"`
	Total   int64             `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	findOneResponse, err := a.store.FindMany(
		&sessionStore.FindManyRequest{
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
	Session session.Session `json:"session"`
}

type UpdateOneResponse struct {
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	if _, err := a.store.UpdateOne(
		&sessionStore.UpdateOneRequest{
			Session: request.Session,
		},
	); err != nil {
		return err
	}
	return nil
}
