package admin

import "time"

type Admin interface {
	XLSXStandardBankStatementToXLSXBudget(*XLSXStandardBankStatementToXLSXBudgetRequest) (*XLSXStandardBankStatementToXLSXBudgetResponse, error)
}

const ServiceProvider = "Budget-Admin"
const XLSXStandardBankStatementToXLSXBudget = ServiceProvider + ".XLSXStandardBankStatementToXLSXBudget"

type XLSXStandardBankStatementToXLSXBudgetRequest struct {
	XLSXStatement []byte
}

type XLSXStandardBankStatementToXLSXBudgetResponse struct {
	XLSXBudgets map[string]map[time.Month][]byte
}
