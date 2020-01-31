package basic

import (
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
)

type admin struct {
}

func New() budgetAdmin.Admin {
	return &admin{}
}

func (a admin) GetBudgetForMonthInYear(*budgetAdmin.GetBudgetForMonthInYearRequest) (*budgetAdmin.GetBudgetForMonthInYearResponse, error) {
	panic("implement me")
}
