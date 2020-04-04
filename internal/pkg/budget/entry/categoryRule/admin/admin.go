package admin

import (
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Admin interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
	CategoriseBudgetEntry(CategoriseBudgetEntryRequest) (*CategoriseBudgetEntryResponse, error)
}

type CreateOneRequest struct {
	Claims       claims.Claims                        `validate:"required"`
	CategoryRule budgetEntryCategoryRule.CategoryRule `validate:"required"`
}

type CreateOneResponse struct {
	CategoryRule budgetEntryCategoryRule.CategoryRule
}

type UpdateOneRequest struct {
	Claims       claims.Claims                        `validate:"required"`
	CategoryRule budgetEntryCategoryRule.CategoryRule `validate:"required"`
}

type UpdateOneResponse struct {
}

type CategoriseBudgetEntryRequest struct {
	Claims                 claims.Claims `validate:"required"`
	BudgetEntryDescription string        `validate:"required"`
}

type CategoriseBudgetEntryResponse struct {
	CategoryRule budgetEntryCategoryRule.CategoryRule
}
