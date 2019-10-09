package jsonRpc

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/authenticator"
	"net/http"
)

type adaptor struct {
	authenticator authenticator.Authenticator
}

func New(
	authenticator authenticator.Authenticator,
) *adaptor {
	return &adaptor{
		authenticator: authenticator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return authenticator.ServiceProvider
}

type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpResponse struct {
	JWT string `json:"jwt"`
}

func (a *adaptor) Login(r *http.Request, request *SignUpRequest, response *SignUpResponse) error {
	loginResponse, err := a.authenticator.Login(
		&authenticator.LoginRequest{
			Email:    request.Email,
			Password: request.Password,
		},
	)
	if err != nil {
		return err
	}

	response.JWT = loginResponse.JWT

	return nil
}
