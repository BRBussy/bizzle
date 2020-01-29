package jsonRpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"net/http"
)

type adaptor struct {
	store roleStore.Store
}

func New(
	store roleStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return roleStore.ServiceProvider
}

type CreateOneRequest struct {
	Role role.Role `json:"role"`
}

type CreateOneResponse struct {
	Role role.Role `json:"role"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		&roleStore.CreateOneRequest{
			Role: request.Role,
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
	Role role.Role `json:"role"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	findOneResponse, err := a.store.FindOne(
		&roleStore.FindOneRequest{
			Identifier: request.Identifier.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.Role = findOneResponse.Role

	return nil
}

type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

type FindManyResponse struct {
	Records []role.Role `json:"records"`
	Total   int64       `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	findOneResponse, err := a.store.FindMany(
		&roleStore.FindManyRequest{
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
	Role role.Role `json:"role"`
}

type UpdateOneResponse struct {
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	if _, err := a.store.UpdateOne(
		&roleStore.UpdateOneRequest{
			Role: request.Role,
		},
	); err != nil {
		return err
	}
	return nil
}
