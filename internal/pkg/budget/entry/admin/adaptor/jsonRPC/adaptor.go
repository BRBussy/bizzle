package jsonRPC

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/rs/zerolog/log"
	"net/http"
)

type adaptor struct {
	admin budgetEntryAdmin.Admin
}

func New(
	admin budgetEntryAdmin.Admin,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		admin: admin,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetEntryAdmin.ServiceProvider
}

type CreateManyRequest struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

type CreateManyResponse struct {
}

func (a *adaptor) CreateMany(r *http.Request, request *CreateManyRequest, response *CreateManyResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err)
		return exception.ErrUnexpected{}
	}

	if _, err := a.admin.CreateMany(&budgetEntryAdmin.CreateManyRequest{
		BudgetEntries: request.BudgetEntries,
		Claims:        c,
	}); err != nil {
		return err
	}

	return nil
}
