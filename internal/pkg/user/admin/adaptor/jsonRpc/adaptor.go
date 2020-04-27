package jsonRpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
	"net/http"
)

type adaptor struct {
	userAdmin userAdmin.Admin
}

func New(
	userAdmin userAdmin.Admin,
) *adaptor {
	return &adaptor{
		userAdmin: userAdmin,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return userAdmin.ServiceProvider
}

type CreateOneRequest struct {
	User user.User `json:"user"`
}

type CreateOneResponse struct {
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return bizzleException.ErrUnexpected{}
	}

	if _, err := a.userAdmin.CreateOne(
		userAdmin.CreateOneRequest{
			Claims: c,
			User:   request.User,
		},
	); err != nil {
		return err
	}

	return nil
}

type RegisterOneRequest struct {
	UserIdentifier identifier.Serialized `json:"userIdentifier"`
	Password       string                `json:"password"`
}

type RegisterOneResponse struct {
}

func (a *adaptor) RegisterOne(r *http.Request, request *RegisterOneRequest, response *RegisterOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return bizzleException.ErrUnexpected{}
	}

	if _, err := a.userAdmin.RegisterOne(
		userAdmin.RegisterOneRequest{
			Claims:         c,
			UserIdentifier: request.UserIdentifier.Identifier,
		},
	); err != nil {
		return err
	}

	return nil
}

type ChangePasswordRequest struct {
	UserIdentifier identifier.Serialized `json:"userIdentifier"`
	Password       string                `json:"password"`
}

type ChangePasswordResponse struct {
}

func (a *adaptor) ChangePassword(r *http.Request, request *ChangePasswordRequest, response *ChangePasswordResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not pass claims for context")
		return bizzleException.ErrUnexpected{}
	}

	if _, err := a.userAdmin.ChangePassword(
		userAdmin.ChangePasswordRequest{
			Claims:         c,
			UserIdentifier: request.UserIdentifier.Identifier,
		},
	); err != nil {
		return err
	}

	return nil
}
