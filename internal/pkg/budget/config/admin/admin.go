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

type GetMyConfigRequest struct {
	Claims claims.Claims `validate:"required"`
}

type GetMyConfigResponse struct {
	Config budgetConfig.Config
}

type SetMyConfigRequest struct {
	Claims claims.Claims `validate:"required"`
	Config budgetConfig.Config
}

type SetMyConfigResponse struct {
}
