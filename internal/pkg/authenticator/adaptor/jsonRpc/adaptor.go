package jsonRpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"net/http"
)

type adaptor struct {
	authenticator bizzleAuthenticator.Authenticator
}

func New(
	authenticator bizzleAuthenticator.Authenticator,
) *adaptor {
	return &adaptor{
		authenticator: authenticator,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return bizzleAuthenticator.ServiceProvider
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
		&bizzleAuthenticator.LoginRequest{
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

type AuthenticateServiceRequest struct {
	Claims  claims.Serialized `json:"claims"`
	Service string            `json:"service"`
}

type AuthenticateServiceResponse struct {
}

func (a *adaptor) AuthenticateService(r *http.Request, request *AuthenticateServiceRequest, response *AuthenticateServiceResponse) error {
	if _, err := a.authenticator.AuthenticateService(
		&bizzleAuthenticator.AuthenticateServiceRequest{
			Claims:  request.Claims.Claims,
			Service: request.Service,
		},
	); err != nil {
		return err
	}
	return nil
}
