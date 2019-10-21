package jsonRPC

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/security/token/generator"
	"net/http"
)

type adaptor struct {
	generator generator.Generator
}

func New(
	generator generator.Generator,
) *adaptor {
	return &adaptor{
		generator: generator,
	}
}

func (a *adaptor) ServiceName() string {
	return "TokenGenerator"
}

type GenerateTokenRequest struct {
	Claims claims.Serialized `json:"claims"`
}

type GenerateTokenResponse struct {
	Token string `json:"token"`
}

func (a *adaptor) GenerateToken(r *http.Request, request *GenerateTokenRequest, response *GenerateTokenResponse) error {
	generateResponse, err := a.generator.GenerateToken(
		r.Context(),
		&generator.GenerateTokenRequest{
			Claims: request.Claims.Claims,
		},
	)
	if err != nil {
		return err
	}

	response.Token = generateResponse.Token
	return nil
}
