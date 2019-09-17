package jsonRpc

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"net/http"
)

type adaptor struct {
	store roleStore.Store
}

func New(
	authenticator roleStore.Store,
) *adaptor {
	return &adaptor{
		store: authenticator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return roleStore.ServiceProvider
}

func (a *adaptor) MethodRequiresAuthorization(method string) bool {
	return false
}

type CreateRequest struct {
	Role role.Role `json:"role"`
}

type CreateResponse struct {
	Role role.Role `json:"role"`
}

func (a *adaptor) Create(r *http.Request, request *CreateRequest, response *CreateResponse) error {
	_, err := a.store.Create(
		&roleStore.CreateRequest{
			Role: request.Role,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
