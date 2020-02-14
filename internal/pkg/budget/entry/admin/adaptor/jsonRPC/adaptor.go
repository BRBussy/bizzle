package jsonRPC

import (
	"encoding/base64"
	"net/http"

	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/rs/zerolog/log"
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
		return bizzleException.ErrUnexpected{}
	}

	if _, err := a.admin.CreateMany(&budgetEntryAdmin.CreateManyRequest{
		BudgetEntries: request.BudgetEntries,
		Claims:        c,
	}); err != nil {
		return err
	}

	return nil
}

type XLSXStandardBankStatementToBudgetEntriesRequest struct {
	XLSXStatement string `json:"xlsxStatement"`
}

type XLSXStandardBankStatementToBudgetEntriesResponse struct {
	BudgetCompositeEntries []budgetEntry.CompositeEntry `json:"budgetCompositeEntries"`
}

func (a *adaptor) XLSXStandardBankStatementToBudgetCompositeEntries(r *http.Request, request *XLSXStandardBankStatementToBudgetEntriesRequest, response *XLSXStandardBankStatementToBudgetEntriesResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err)
		return bizzleException.ErrUnexpected{}
	}

	// parse xlsx statement to bytes
	statementFileDataBytes, err := base64.StdEncoding.DecodeString(request.XLSXStatement)
	if err != nil {
		log.Error().Err(err).Msg("failed to decode excel file data")
		return bizzleException.ErrUnexpected{}
	}

	xlsxStandardBankStatementToBudgetEntriesResponse, err := a.admin.XLSXStandardBankStatementToBudgetEntries(&budgetEntryAdmin.XLSXStandardBankStatementToBudgetEntriesRequest{
		XLSXStatement: statementFileDataBytes,
		Claims:        c,
	})
	if err != nil {
		return err
	}

	response.BudgetCompositeEntries = xlsxStandardBankStatementToBudgetEntriesResponse.BudgetEntries

	return nil
}

type DuplicateCheckRequest struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

type DuplicateCheckResponse struct {
	Uniques             []budgetEntry.Entry `json:"uniques"`
	ExactDuplicates     []budgetEntry.Entry `json:"exactDuplicates"`
	SuspectedDuplicates []budgetEntry.Entry `json:"suspectedDuplicates"`
}

func (a *adaptor) DuplicateCheck(r *http.Request, request *DuplicateCheckRequest, response *DuplicateCheckResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	duplicateCheckResponse, err := a.admin.DuplicateCheck(&budgetEntryAdmin.DuplicateCheckRequest{
		BudgetEntries: request.BudgetEntries,
		Claims:        c,
	})
	if err != nil {
		return err
	}
	response.Uniques = duplicateCheckResponse.Uniques
	response.ExactDuplicates = duplicateCheckResponse.ExactDuplicates
	response.SuspectedDuplicates = duplicateCheckResponse.SuspectedDuplicates

	return nil
}
