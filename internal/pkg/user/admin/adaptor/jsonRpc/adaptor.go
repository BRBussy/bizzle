package jsonRpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	"net/http"
)

type adaptor struct {
	admin userAdmin.Admin
}

func New(
	authenticator userAdmin.Admin,
) *adaptor {
	return &adaptor{
		admin: authenticator,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return userAdmin.ServiceProvider
}

func (a *adaptor) MethodRequiresAuthorization(method string) bool {
	return false
}

type CreateOneRequest struct {
	User user.User `json:"user"`
}

type CreateOneResponse struct {
	User user.User `json:"user"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	_, err := a.admin.CreateOne(
		&userAdmin.CreateOneRequest{
			User: request.User,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
