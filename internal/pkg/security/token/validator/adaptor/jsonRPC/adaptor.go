package jsonRPC

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
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
	Claims claims.Serialized `json:"claims"`
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
	response.Claims.Claims = validateResponse.Claims
	return nil
}
