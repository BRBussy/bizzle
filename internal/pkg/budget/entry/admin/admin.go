package admin

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryIgnored "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

// Admin is the budget entry admin interface
type Admin interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	CreateMany(CreateManyRequest) (*CreateManyResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
	UpdateMany(UpdateManyRequest) (*UpdateManyResponse, error)
	DuplicateCheck(DuplicateCheckRequest) (*DuplicateCheckResponse, error)
	XLSXStandardBankStatementToBudgetEntries(XLSXStandardBankStatementToBudgetEntriesRequest) (*XLSXStandardBankStatementToBudgetEntriesResponse, error)
	DeleteOne(DeleteOneRequest) (*DeleteOneResponse, error)
	IgnoreOne(IgnoreOneRequest) (*IgnoreOneResponse, error)
	RecogniseOne(RecogniseOneRequest) (*RecogniseOneResponse, error)
	IgnoredCheck(IgnoredCheckRequest) (*IgnoredCheckResponse, error)
}

// ServiceProvider is the budget entry admin service provider name
const ServiceProvider = "BudgetEntry-Admin"

// CreateManyService is the service name for CreateMany
const CreateManyService = ServiceProvider + ".CreateMany"

// CreateOneService is the service name for CreateMany
const CreateOneService = ServiceProvider + ".CreateOne"

// UpdateOneService is the service name for UpdateOne
const UpdateOneService = ServiceProvider + ".UpdateOne"

// UpdateManyService is the service name for UpdateMany
const UpdateManyService = ServiceProvider + ".UpdateMany"

// DuplicateCheckService is the service name for DuplicateCheck
const DuplicateCheckService = ServiceProvider + ".DuplicateCheck"

// XLSXStandardBankStatementToBudgetCompositeEntriesService is the service name for XLSXStandardBankStatementToBudgetCompositeEntries
const XLSXStandardBankStatementToBudgetCompositeEntriesService = ServiceProvider + ".XLSXStandardBankStatementToBudgetEntries"

// DeleteOneService is the service name for DeleteOne
const DeleteOneService = ServiceProvider + ".DeleteOne"

// IgnoreOneService is the service name for DeleteOne
const IgnoreOneService = ServiceProvider + ".IgnoreOne"

// RecogniseOneService is the service name for DeleteOne
const RecogniseOneService = ServiceProvider + ".RecogniseOne"

// IgnoredCheckService is the service name for DeleteOne
const IgnoredCheckService = ServiceProvider + ".IgnoredCheck"

// CreateManyRequest is the request object for the CreateMany service
type CreateManyRequest struct {
	BudgetEntries []budgetEntry.Entry `validate:"required,gt=0"`
	Claims        claims.Claims       `validate:"required"`
}

// CreateManyResponse is the response object for the CreateMany service
type CreateManyResponse struct {
}

// DuplicateCheckRequest is the request object for the DuplicateCheck service
type DuplicateCheckRequest struct {
	BudgetEntries []budgetEntry.Entry `validate:"required,gt=1"`
	Claims        claims.Claims       `validate:"required"`
}

// DuplicateCheckResponse is the response object for the DuplicateCheck service
type DuplicateCheckResponse struct {
	Uniques             []budgetEntry.Entry
	ExactDuplicates     []DuplicateEntries
	SuspectedDuplicates []DuplicateEntries
}

// DuplicateEntries is used to hold an duplicate check result pair
type DuplicateEntries struct {
	Existing budgetEntry.Entry `json:"existing"`
	New      budgetEntry.Entry `json:"new"`
}

// XLSXStandardBankStatementToBudgetEntriesRequest is the XLSXStandardBankStatementToBudgetEntries service request object
type XLSXStandardBankStatementToBudgetEntriesRequest struct {
	Claims        claims.Claims `validate:"required"`
	XLSXStatement []byte        `validate:"required"`
}

// XLSXStandardBankStatementToBudgetEntriesResponse is the XLSXStandardBankStatementToBudgetEntries service response object
type XLSXStandardBankStatementToBudgetEntriesResponse struct {
	BudgetEntries []budgetEntry.Entry
}

// UpdateOneRequest is the request object for the UpdateOne service
type UpdateOneRequest struct {
	Claims      claims.Claims `validate:"required"`
	BudgetEntry budgetEntry.Entry
}

// UpdateOneResponse is the response object for the UpdateOneService
type UpdateOneResponse struct {
}

// CreateOneRequest is the request object for the CreateOne service
type CreateOneRequest struct {
	Claims      claims.Claims `validate:"required"`
	BudgetEntry budgetEntry.Entry
}

// CreateOneResponse is the response object for the CreateOneService
type CreateOneResponse struct {
	BudgetEntry budgetEntry.Entry
}

// UpdateManyRequest is the request object for the UpdateMany service
type UpdateManyRequest struct {
	Claims        claims.Claims `validate:"required"`
	BudgetEntries []budgetEntry.Entry
}

// UpdateManyResponse is the response object for the UpdateMany service
type UpdateManyResponse struct {
}

// DeleteOneRequest is the request object for the DeleteOne service
type DeleteOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

// DeleteOneResponse is the response object for the the DeleteOne service
type DeleteOneResponse struct {
}

type IgnoreOneRequest struct {
	Claims      claims.Claims     `validate:"required"`
	BudgetEntry budgetEntry.Entry `validate:"required"`
}

type IgnoreOneResponse struct {
}

type RecogniseOneRequest struct {
	Claims      claims.Claims     `validate:"required"`
	BudgetEntry budgetEntry.Entry `validate:"required"`
}

type RecogniseOneResponse struct {
}

type IgnoredCheckRequest struct {
	Claims        claims.Claims       `validate:"required"`
	BudgetEntries []budgetEntry.Entry `validate:"required"`
}

type IgnoredCheckResponse struct {
	Ignored []budgetEntryIgnored.Ignored
}
