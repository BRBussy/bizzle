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
	Config budgetConfig.Config `json:"config"`
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

	response.Config = getMyConfigResponse.Config

	return nil
}
