package jsonRpc

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/exercise/session"
	sessionAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/session/admin"
	"net/http"
)

type adaptor struct {
	store sessionAdmin.Admin
}

func New(
	store sessionAdmin.Admin,
) jsonRpcServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return sessionAdmin.ServiceProvider
}

type CreateOneRequest struct {
	Session session.Session `json:"session"`
}

type CreateOneResponse struct {
	Session session.Session `json:"session"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	if _, err := a.store.CreateOne(
		&sessionAdmin.CreateOneRequest{
			Session: request.Session,
		},
	); err != nil {
		return err
	}

	return nil
}
