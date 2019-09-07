package jsonRpc

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/package/authenticator"
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

func (a *adaptor) MethodRequiresAuthorization(method string) bool {
	return false
}

type SignUpRequest struct {
	Username string `json:"username"`
}

type SignUpResponse struct {
	Msg string `json:"msg"`
}

func (a *adaptor) SignUp(r *http.Request, request *SignUpRequest, response *SignUpResponse) error {
	signUpResponse, err := a.authenticator.SignUp(
		&authenticator.SignUpRequest{},
	)
	if err != nil {
		return err
	}

	response.Msg = signUpResponse.Msg

	return nil
}
