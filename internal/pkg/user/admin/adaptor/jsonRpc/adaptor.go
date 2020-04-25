package jsonRpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
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
	User user.User `json:"user"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.userAdmin.CreateOne(
		userAdmin.CreateOneRequest{
			User: request.User,
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
}
