package admin

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Admin interface {
	CreateMany(*CreateManyRequest) (*CreateManyResponse, error)
	DuplicateCheck(*DuplicateCheckRequest) (*DuplicateCheckResponse, error)
	XLSXStandardBankStatementToBudgetEntries(*XLSXStandardBankStatementToBudgetEntriesRequest) (*XLSXStandardBankStatementToBudgetEntriesResponse, error)
}

const ServiceProvider jsonRPCServiceProvider.Name = "BudgetEntry-Admin"

const CreateManyService = ServiceProvider + ".CreateMany"
const DuplicateCheckService = ServiceProvider + ".DuplicateCheck"
const XLSXStandardBankStatementToBudgetEntriesService = ServiceProvider + ".XLSXStandardBankStatementToBudgetEntries"

type CreateManyRequest struct {
	BudgetEntries []budgetEntry.Entry
	Claims        claims.Claims `validate:"required"`
}

type CreateManyResponse struct {
}

type DuplicateCheckRequest struct {
	BudgetEntries []budgetEntry.Entry
	Claims        claims.Claims `validate:"required"`
}

type DuplicateCheckResponse struct {
	Uniques             []budgetEntry.Entry
	ExactDuplicates     []budgetEntry.Entry
	SuspectedDuplicates []budgetEntry.Entry
}

type XLSXStandardBankStatementToBudgetEntriesRequest struct {
	XLSXStatement []byte
}

type XLSXStandardBankStatementToBudgetEntriesResponse struct {
	BudgetEntries []budgetEntry.Entry
}
