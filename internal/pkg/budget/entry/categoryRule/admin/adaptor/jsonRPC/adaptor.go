package jsonRPC

import (
	"net/http"

	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	budgetEntryCategoryRuleAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/rs/zerolog/log"
)

type adaptor struct {
	admin budgetEntryCategoryRuleAdmin.Admin
}

// New creates a new jsonrpc adaptor for a budget entry admin
func New(
	admin budgetEntryCategoryRuleAdmin.Admin,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		admin: admin,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetEntryCategoryRuleAdmin.ServiceProvider
}

// UpdateOneRequest is the request object for UpdateOne method
type UpdateOneRequest struct {
	CategoryRule budgetEntryCategoryRule.CategoryRule `json:"categoryRule"`
}

// UpdateOneResponse is the response object for the update one method
type UpdateOneResponse struct {
}

func (a *adaptor) UpdateOne(r *http.Request, request *UpdateOneRequest, response *UpdateOneResponse) error {
	c, err := claims.ParseClaimsFromContext(r.Context())
	if err != nil {
		log.Error().Err(err).Msg("could not parse claims from context")
		return exception.ErrUnexpected{}
	}

	if _, err := a.admin.UpdateOne(
		budgetEntryCategoryRuleAdmin.UpdateOneRequest{
			Claims:       c,
			CategoryRule: request.CategoryRule,
		},
	); err != nil {
		return err
	}

	return nil
}
