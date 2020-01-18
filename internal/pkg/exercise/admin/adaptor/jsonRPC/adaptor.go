package jsonRPC

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
	exerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin"
	"net/http"
)

type adaptor struct {
	admin exerciseAdmin.Admin
}

func New(
	admin exerciseAdmin.Admin,
) jsonRpcServiceProvider.Provider {
	return &adaptor{
		admin: admin,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return exerciseAdmin.ServiceProvider
}

type CreateOneRequest struct {
	Exercise exercise.Exercise `json:"exercise"`
}

type CreateOneResponse struct {
	Exercise exercise.Exercise `json:"exercise"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	createOneResponse, err := a.admin.CreateOne(
		&exerciseAdmin.CreateOneRequest{
			Exercise: request.Exercise,
		},
	)
	if err != nil {
		return err
	}

	response.Exercise = createOneResponse.Exercise

	return nil
}
