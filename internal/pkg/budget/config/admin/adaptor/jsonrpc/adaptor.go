package jsonrpc

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetConfig "github.com/BRBussy/bizzle/internal/pkg/budget/config"
	budgetConfigAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/config/admin"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/rs/zerolog/log"
	"net/http"
)

type adaptor struct {
	admin budgetConfigAdmin.Admin
}

// New creates a new jsonrpc adaptor for a budget config admin
func New(
	admin budgetConfigAdmin.Admin,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		admin: admin,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetConfigAdmin.ServiceProvider
}

type GetMyConfigRequest struct {
}

type GetMyConfigResponse struct {
	BudgetConfig budgetConfig.Config `json:"budgetConfig"`
}

func (a *adaptor) GetMyConfig(r *http.Request, request *GetMyConfigRequest, response *GetMyConfigResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err)
		return bizzleException.ErrUnexpected{}
	}

	getMyConfigResponse, err := a.admin.GetMyConfig(
		budgetConfigAdmin.GetMyConfigRequest{
			Claims: c,
		},
	)
	if err != nil {
		return err
	}

	response.BudgetConfig = getMyConfigResponse.BudgetConfig

	return nil
}

type SetMyConfigRequest struct {
	BudgetConfig budgetConfig.Config `json:"budgetConfig"`
}

type SetMyConfigResponse struct {
}

func (a *adaptor) SetMyConfig(r *http.Request, request *SetMyConfigRequest, response *SetMyConfigResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err)
		return bizzleException.ErrUnexpected{}
	}

	if _, err := a.admin.SetMyConfig(
		budgetConfigAdmin.SetMyConfigRequest{
			Claims:       c,
			BudgetConfig: request.BudgetConfig,
		},
	); err != nil {
		return err
	}

	return nil
}
