package jsonRPC

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type adaptor struct {
	admin budgetAdmin.Admin
}

func New(
	admin budgetAdmin.Admin,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		admin: admin,
	}
}

func (a adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetAdmin.ServiceProvider
}

type GetBudgetForDateRangeRequest struct {
	StartDate time.Time `json:"startDate"`
	EndDate   time.Time `json:"endDate"`
}

type GetBudgetForDateRangeResponse struct {
	Budget budget.Budget `json:"budget"`
}

func (a *adaptor) GetBudgetForDateRange(r *http.Request, request *GetBudgetForDateRangeRequest, response *GetBudgetForDateRangeResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err)
		return bizzleException.ErrUnexpected{}
	}

	getBudgetForMonthInYearResponse, err := a.admin.GetBudgetForDateRange(&budgetAdmin.GetBudgetForDateRangeRequest{
		Claims:    c,
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
	})
	if err != nil {
		return err
	}

	response.Budget = getBudgetForMonthInYearResponse.Budget

	return nil
}
