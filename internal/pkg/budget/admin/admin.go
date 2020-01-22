package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	"time"
)

type Admin interface {
	XLSXStandardBankStatementToXLSXBudget(*XLSXStandardBankStatementToXLSXBudgetRequest) (*XLSXStandardBankStatementToXLSXBudgetResponse, error)
	BudgetEntriesToBudgets(*BudgetEntriesToBudgetsRequest) (*BudgetEntriesToBudgetResponse, error)
	BudgetToXLSX(*BudgetToXLSXRequest) (*BudgetToXLSXResponse, error)
}

const ServiceProvider = "Budget-Admin"
const XLSXStandardBankStatementToXLSXBudget = ServiceProvider + ".XLSXStandardBankStatementToXLSXBudget"

type XLSXStandardBankStatementToXLSXBudgetRequest struct {
	XLSXStatement []byte
}

type XLSXStandardBankStatementToXLSXBudgetResponse struct {
	XLSXBudgets map[string]map[time.Month][]byte
}

type BudgetEntriesToBudgetsRequest struct {
	BudgetEntries []budget.Entry
}

type BudgetEntriesToBudgetResponse struct {
	Budgets []budget.Budget
}

type BudgetToXLSXRequest struct {
	Budget budget.Budget
}

type BudgetToXLSXResponse struct {
	XLSXBudget []byte
}
