package jsonRPC

import (
	"encoding/base64"
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	"github.com/rs/zerolog/log"
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
	XLSXStatement string `json:"xlsxStatement"`
}

type XLSXStandardBankStatementToXLSXBudgetResponse struct {
	XLSXBudgets map[string]map[time.Month][]byte `json:"xlsxBudgets"`
}

func (a *adaptor) XLSXStandardBankStatementToXLSXBudget(r *http.Request, request *XLSXStandardBankStatementToXLSXBudgetRequest, response *XLSXStandardBankStatementToXLSXBudgetResponse) error {
	// parse xlsx statement to bytes
	statementFileDataBytes, err := base64.StdEncoding.DecodeString(request.XLSXStatement)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode picture data")
		return err
	}

	// call service
	processResponse, err := a.admin.XLSXStandardBankStatementToXLSXBudget(
		&budgetAdmin.XLSXStandardBankStatementToXLSXBudgetRequest{
			XLSXStatement: statementFileDataBytes,
		},
	)
	if err != nil {
		return err
	}

	response.XLSXBudgets = processResponse.XLSXBudgets

	return nil
}
