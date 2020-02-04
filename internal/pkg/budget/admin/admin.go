package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"time"
)

type Admin interface {
	GetBudgetForMonthInYear(*GetBudgetForMonthInYearRequest) (*GetBudgetForMonthInYearResponse, error)
	GetBudgetForDateRange(*GetBudgetForDateRangeRequest) (*GetBudgetForDateRangeResponse, error)
}

const ServiceProvider = "Budget-Admin"
const GetBudgetForMonthInYearService = ServiceProvider + ".GetBudgetForMonthInYear"
const GetBudgetForDateRangeService = ServiceProvider + ".GetBudgetForDateRange"

type GetBudgetForMonthInYearRequest struct {
	Claims claims.Claims `validate:"required"`
	Month  time.Month    `validate:"required"`
}

type GetBudgetForMonthInYearResponse struct {
	Budget budget.Budget
}

type GetBudgetForDateRangeRequest struct {
	Claims    claims.Claims `validate:"required"`
	StartDate time.Time     `validate:"required"`
	EndDate   time.Time     `validate:"required"`
}

type GetBudgetForDateRangeResponse struct {
	Budget budget.Budget
}
