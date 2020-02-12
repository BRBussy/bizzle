package jsonRPC

import (
	jsonRPCServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	budgetCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type adaptor struct {
	store budgetCategoryRuleStore.Store
}

// New creates a new budget category rule store json rpc adaptor
func New(
	store budgetCategoryRuleStore.Store,
) jsonRPCServiceProvider.Provider {
	return &adaptor{
		store: store,
	}
}

func (a *adaptor) Name() jsonRPCServiceProvider.Name {
	return budgetCategoryRuleStore.ServiceProvider
}

type FindOneRequest struct {
	Identifier identifier.Serialized `json:"identifier"`
}

type FindOneResponse struct {
	BudgetCategoryRule budgetCategoryRule.CategoryRule `json:"categoryRule"`
}
