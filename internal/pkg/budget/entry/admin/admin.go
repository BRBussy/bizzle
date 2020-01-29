package admin

import (
	jsonRpcServiceProvider "github.com/BRBussy/bizzle/internal/pkg/api/jsonRpc/service/provider"
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
)

type Admin interface {
	CreateMany(*CreateManyRequest) (*CreateManyResponse, error)
}

const ServiceProvider jsonRpcServiceProvider.Name = "BudgetEntry-Admin"

type CreateManyRequest struct {
	BudgetEntries []budgetEntry.Entry
}

type CreateManyResponse struct {
}
