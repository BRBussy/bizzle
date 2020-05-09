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
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

type adaptor struct {
	admin budgetEntryAdmin.Admin
}

// New creates a new jsonrpc adaptor for a budget entry admin
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

// CreateManyRequest is the request object for CreateMany method
type CreateManyRequest struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

// CreateManyResponse is the response object for the CreateMany method
type CreateManyResponse struct {
}

func (a *adaptor) CreateMany(r *http.Request, request *CreateManyRequest, response *CreateManyResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err)
		return bizzleException.ErrUnexpected{}
	}

	if _, err := a.admin.CreateMany(budgetEntryAdmin.CreateManyRequest{
		BudgetEntries: request.BudgetEntries,
		Claims:        c,
	}); err != nil {
		return err
	}

	return nil
}

// XLSXStandardBankStatementToBudgetEntriesRequest is the request object for XLSXStandardBankStatementToBudgetEntries method
type XLSXStandardBankStatementToBudgetEntriesRequest struct {
	XLSXStatement string `json:"xlsxStatement"`
}

// XLSXStandardBankStatementToBudgetEntriesResponse is the response object for the XLSXStandardBankStatementToBudgetEntries method
type XLSXStandardBankStatementToBudgetEntriesResponse struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

func (a *adaptor) XLSXStandardBankStatementToBudgetEntries(r *http.Request, request *XLSXStandardBankStatementToBudgetEntriesRequest, response *XLSXStandardBankStatementToBudgetEntriesResponse) error {
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

	xlsxStandardBankStatementToBudgetEntriesResponse, err := a.admin.XLSXStandardBankStatementToBudgetEntries(budgetEntryAdmin.XLSXStandardBankStatementToBudgetEntriesRequest{
		XLSXStatement: statementFileDataBytes,
		Claims:        c,
	})
	if err != nil {
		return err
	}

	response.BudgetEntries = xlsxStandardBankStatementToBudgetEntriesResponse.BudgetEntries

	return nil
}

// DuplicateCheckRequest is the request object for the DuplicateCheck method
type DuplicateCheckRequest struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

// DuplicateCheckResponse is the response object for the DuplicateCheck method
type DuplicateCheckResponse struct {
	Uniques             []budgetEntry.Entry                 `json:"uniques"`
	ExactDuplicates     []budgetEntryAdmin.DuplicateEntries `json:"exactDuplicates"`
	SuspectedDuplicates []budgetEntryAdmin.DuplicateEntries `json:"suspectedDuplicates"`
}

func (a *adaptor) DuplicateCheck(r *http.Request, request *DuplicateCheckRequest, response *DuplicateCheckResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	duplicateCheckResponse, err := a.admin.DuplicateCheck(budgetEntryAdmin.DuplicateCheckRequest{
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

// UpdateOneRequest is the request object for UpdateOne method
type UpdateOneRequest struct {
	BudgetEntry budgetEntry.Entry `json:"budgetEntry"`
}

// UpdateOneResponse is the response object for the update one method
type UpdateOneResponse struct {
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	if _, err := a.admin.UpdateOne(budgetEntryAdmin.UpdateOneRequest{
		Claims:      c,
		BudgetEntry: request.BudgetEntry,
	}); err != nil {
		return err
	}

	return nil
}

// CreateOneRequest is the request object for CreateOne method
type CreateOneRequest struct {
	BudgetEntry budgetEntry.Entry `json:"budgetEntry"`
}

// CreateOneResponse is the response object for the CreateOne method
type CreateOneResponse struct {
	BudgetEntry budgetEntry.Entry `json:"budgetEntry"`
}

func (a *adaptor) CreateOne(r *http.Request, request *CreateOneRequest, response *CreateOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	createOneResponse, err := a.admin.CreateOne(budgetEntryAdmin.CreateOneRequest{
		Claims:      c,
		BudgetEntry: request.BudgetEntry,
	})
	if err != nil {
		return err
	}

	response.BudgetEntry = createOneResponse.BudgetEntry

	return nil
}

// UpdateManyRequest is the request object for UpdateMany method
type UpdateManyRequest struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

// UpdateManyResponse is the response object for the UpdateMany method
type UpdateManyResponse struct {
}

func (a *adaptor) UpdateMany(r *http.Request, request *UpdateManyRequest, response *UpdateManyResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	if _, err := a.admin.UpdateMany(budgetEntryAdmin.UpdateManyRequest{
		Claims:        c,
		BudgetEntries: request.BudgetEntries,
	}); err != nil {
		return err
	}

	return nil
}

// DeleteOneRequest is the request object for DeleteOne method
type DeleteOneRequest struct {
	Identifier identifier.Serialized `json:"identifier"`
}

// DeleteOneResponse is the response object for the update one method
type DeleteOneResponse struct {
}

func (a *adaptor) DeleteOne(r *http.Request, request *DeleteOneRequest, response *DeleteOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	if _, err := a.admin.DeleteOne(budgetEntryAdmin.DeleteOneRequest{
		Claims:     c,
		Identifier: request.Identifier.Identifier,
	}); err != nil {
		return err
	}

	return nil
}

type IgnoreOneRequest struct {
	Description string `json:"description"`
}

type IgnoreOneResponse struct {
}

func (a *adaptor) IgnoreOne(r *http.Request, request *IgnoreOneRequest, response *IgnoreOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	if _, err := a.admin.IgnoreOne(
		budgetEntryAdmin.IgnoreOneRequest{
			Claims:      c,
			Description: request.Description,
		},
	); err != nil {
		return err
	}

	return nil
}

type IgnoredCheckRequest struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

type IgnoredCheckResponse struct {
	BudgetEntries []budgetEntry.Entry `json:"budgetEntries"`
}

func (a *adaptor) IgnoredCheck(r *http.Request, request *IgnoredCheckRequest, response *IgnoredCheckResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	ignoredCheckResponse, err := a.admin.IgnoredCheck(
		budgetEntryAdmin.IgnoredCheckRequest{
			Claims:        c,
			BudgetEntries: request.BudgetEntries,
		},
	)
	if err != nil {
		return err
	}

	response.BudgetEntries = ignoredCheckResponse.BudgetEntries

	return nil
}
