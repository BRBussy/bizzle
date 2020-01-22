package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	"time"
)

type Admin interface {
	XLSXStandardBankStatementToXLSXBudget(*XLSXStandardBankStatementToXLSXBudgetRequest) (*XLSXStandardBankStatementToXLSXBudgetResponse, error)
	BudgetEntriesToXLSXBudget(*BudgetEntriesToXLSXBudgetRequest) (*BudgetEntriesToXLSXBudgetResponse, error)
}

const ServiceProvider = "Budget-Admin"
const XLSXStandardBankStatementToXLSXBudget = ServiceProvider + ".XLSXStandardBankStatementToXLSXBudget"

type XLSXStandardBankStatementToXLSXBudgetRequest struct {
	XLSXStatement []byte
}

type XLSXStandardBankStatementToXLSXBudgetResponse struct {
	XLSXBudgets map[string]map[time.Month][]byte
}

type BudgetEntriesToXLSXBudgetRequest struct {
	BudgetEntries []budget.Entry
}

type BudgetEntriesToXLSXBudgetResponse struct {
	XLSXBudget []byte
}
