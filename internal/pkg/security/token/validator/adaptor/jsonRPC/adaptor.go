package jsonRPC

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	tokenValidator "github.com/BRBussy/bizzle/internal/pkg/security/token/validator"
	"net/http"
)

type adaptor struct {
	validator tokenValidator.Validator
}

func New(validator tokenValidator.Validator) jsonRpcServiceProvider.Provider {
	return &adaptor{
		validator: validator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return tokenValidator.ServiceProvider
}

type ValidateRequest struct {
	Token string `json:"token"`
}

type ValidateResponse struct {
	MarshalledClaims []byte `json:"marshalledClaims"`
}

func (a *adaptor) Validate(r *http.Request, request *ValidateRequest, response *ValidateResponse) error {
	validateResponse, err := a.validator.Validate(
		&tokenValidator.ValidateRequest{
			Token: request.Token,
		},
	)
	if err != nil {
		return err
	}
	response.MarshalledClaims = validateResponse.MarshalledClaims
	return nil
}
