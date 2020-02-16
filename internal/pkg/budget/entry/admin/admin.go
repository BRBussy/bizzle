package admin

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Admin interface {
	CreateMany(*CreateManyRequest) (*CreateManyResponse, error)
	DuplicateCheck(*DuplicateCheckRequest) (*DuplicateCheckResponse, error)
	XLSXStandardBankStatementToBudgetEntries(*XLSXStandardBankStatementToBudgetEntriesRequest) (*XLSXStandardBankStatementToBudgetEntriesResponse, error)
}

const ServiceProvider = "BudgetEntry-Admin"

const CreateManyService = ServiceProvider + ".CreateMany"
const DuplicateCheckService = ServiceProvider + ".DuplicateCheck"
const XLSXStandardBankStatementToBudgetCompositeEntriesService = ServiceProvider + ".XLSXStandardBankStatementToBudgetEntries"

type CreateManyRequest struct {
	BudgetEntries []budgetEntry.Entry `validate:"required,gt=1"`
	Claims        claims.Claims       `validate:"required"`
}

type CreateManyResponse struct {
}

type DuplicateCheckRequest struct {
	BudgetEntries []budgetEntry.Entry `validate:"required,gt=1"`
	Claims        claims.Claims       `validate:"required"`
}

type DuplicateCheckResponse struct {
	Uniques             []budgetEntry.Entry
	ExactDuplicates     []DuplicateEntries
	SuspectedDuplicates []DuplicateEntries
}

type DuplicateEntries struct {
	Existing budgetEntry.Entry `json:"existing"`
	New      budgetEntry.Entry `json:"new"`
}

type XLSXStandardBankStatementToBudgetEntriesRequest struct {
	Claims        claims.Claims `validate:"required"`
	XLSXStatement []byte        `validate:"required"`
}

type XLSXStandardBankStatementToBudgetEntriesResponse struct {
	BudgetEntries []budgetEntry.Entry
}
