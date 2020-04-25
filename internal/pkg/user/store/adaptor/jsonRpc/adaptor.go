package jsonRpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
	"net/http"
)

type adaptor struct {
	store userStore.Store
}

func New(
	authenticator userStore.Store,
) *adaptor {
	return &adaptor{
		store: authenticator,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return userStore.ServiceProvider
}

type CreateOneRequest struct {
	User user.User `json:"user"`
}

type CreateOneResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		userStore.CreateOneRequest{
			User: request.User,
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
	User user.User `json:"user"`
}

func (a *adaptor) FindOne(r *http.Request, request *FindOneRequest, response *FindOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return bizzleException.ErrUnexpected{}
	}

	findOneResponse, err := a.store.FindOne(
		userStore.FindOneRequest{
			Claims:     c,
			Identifier: request.Identifier.Identifier,
		},
	)
	if err != nil {
		return err
	}

	response.User = findOneResponse.User

	return nil
}

type FindManyRequest struct {
	Criteria criteria.Serialized `json:"criteria"`
	Query    mongo.Query         `json:"query"`
}

type FindManyResponse struct {
	Records []user.User `json:"records"`
	Total   int64       `json:"total"`
}

func (a *adaptor) FindMany(r *http.Request, request *FindManyRequest, response *FindManyResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return bizzleException.ErrUnexpected{}
	}

	findOneResponse, err := a.store.FindMany(
		userStore.FindManyRequest{
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
	User user.User `json:"user"`
}

type UpdateOneResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return bizzleException.ErrUnexpected{}
	}

	if _, err := a.store.UpdateOne(
		userStore.UpdateOneRequest{
			Claims: c,
			User:   request.User,
		},
	); err != nil {
		return err
	}

	return nil
}
