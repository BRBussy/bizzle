package jsonRPC

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	"net/http"
	"time"
)

type adaptor struct {
	admin budgetAdmin.Admin
}

func New(
	admin budgetAdmin.Admin,
) jsonRpcServiceProvider.Provider {
	return &adaptor{
		admin: admin,
	}
}

func (a *adaptor) Name() jsonRpcServiceProvider.Name {
	return budgetAdmin.ServiceProvider
}

type XLSXStandardBankStatementToXLSXBudgetRequest struct {
	XLSXStatement []byte `json:"xlsxStatement"`
}

type XLSXStandardBankStatementToXLSXBudgetResponse struct {
	XLSXBudgets map[string]map[time.Month][]byte `json:"xlsxBudgets"`
}

func (a *adaptor) XLSXStandardBankStatementToXLSXBudget(r *http.Request, request *XLSXStandardBankStatementToXLSXBudgetRequest, response *XLSXStandardBankStatementToXLSXBudgetResponse) error {
	processResponse, err := a.admin.XLSXStandardBankStatementToXLSXBudget(
		&budgetAdmin.XLSXStandardBankStatementToXLSXBudgetRequest{
			XLSXStatement: request.XLSXStatement,
		},
	)
	if err != nil {
		return err
	}

	response.XLSXBudgets = processResponse.XLSXBudgets

	return nil
}
