package jsonRpc

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"net/http"
)

type adaptor struct {
	store exerciseStore.Store
}

func New(
	authenticator exerciseStore.Store,
) *adaptor {
	return &adaptor{
		store: authenticator,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return exerciseStore.ServiceProvider
}

func (a *adaptor) MethodRequiresAuthorization(method string) bool {
	return false
}

type FindRequest struct {
	SerializedCriteria criteria.SerializedCriteria `json:"criteria"`
}

type FindResponse struct {
}

func (a *adaptor) Find(r *http.Request, request *FindRequest, response *FindResponse) error {
	_, err := a.store.Find(
		&exerciseStore.FindRequest{
			Criteria: request.SerializedCriteria.Criteria,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
