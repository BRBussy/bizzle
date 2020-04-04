package admin

import (
	budgetConfig "github.com/BRBussy/bizzle/internal/pkg/budget/config"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

// Admin provides services to administer a user's budget config
type Admin interface {
	GetMyConfig(GetMyConfigRequest) (*GetMyConfigResponse, error)
	SetMyConfig(SetMyConfigRequest) (*SetMyConfigResponse, error)
}

const ServiceProvider = "BudgetConfig-Store"

const GetMyConfigService = ServiceProvider + ".GetMyConfig"
const SetMyConfigService = ServiceProvider + ".SetMyConfig"

type GetMyConfigRequest struct {
	Claims claims.Claims `validate:"required"`
}

type GetMyConfigResponse struct {
	BudgetConfig budgetConfig.Config
}

type SetMyConfigRequest struct {
	Claims       claims.Claims `validate:"required"`
	BudgetConfig budgetConfig.Config
}

type SetMyConfigResponse struct {
}
