package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"time"
)

type Admin interface {
	GetBudgetForMonthInYear(*GetBudgetForMonthInYearRequest) (*GetBudgetForMonthInYearResponse, error)
}

type GetBudgetForMonthInYearRequest struct {
	Claims claims.Claims `validate:"required"`
	Month  time.Month    `validate:"required"`
}

type GetBudgetForMonthInYearResponse struct {
	Budget budget.Budget
}
