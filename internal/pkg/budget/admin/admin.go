package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	"time"
)

type Admin interface {
	XLSXStandardBankStatementToXLSXBudgets(*XLSXStandardBankStatementToXLSXBudgetsRequest) (*XLSXStandardBankStatementToXLSXBudgetsResponse, error)
	BudgetEntriesToBudgets(*BudgetEntriesToBudgetsRequest) (*BudgetEntriesToBudgetResponse, error)
	BudgetToXLSX(*BudgetToXLSXRequest) (*BudgetToXLSXResponse, error)
	XLSXStandardBankStatementToBudgets(*XLSXStandardBankStatementToBudgetsRequest) (*XLSXStandardBankStatementToBudgetsResponse, error)
}

const ServiceProvider = "Budget-Admin"
const XLSXStandardBankStatementToXLSXBudget = ServiceProvider + ".XLSXStandardBankStatementToXLSXBudgets"
const XLSXStandardBankStatementBudgets = ServiceProvider + ".XLSXStandardBankStatementToBudgets"

type XLSXStandardBankStatementToXLSXBudgetsRequest struct {
	XLSXStatement []byte
}

type XLSXStandardBankStatementToXLSXBudgetsResponse struct {
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

type XLSXStandardBankStatementToBudgetsRequest struct {
	XLSXStatement []byte
}

type XLSXStandardBankStatementToBudgetsResponse struct {
	Budgets []budget.Budget
}
