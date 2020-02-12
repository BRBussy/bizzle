package parser

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Parser interface {
	ParseStatementToBudgetCompositeEntries(*ParseStatementToBudgetCompositeEntriesRequest) (*ParseStatementToBudgetCompositeEntriesResponse, error)
}

type ParseStatementToBudgetCompositeEntriesRequest struct {
	Claims    claims.Claims `validate:"required"`
	Statement []byte        `validate:"required"`
}

type ParseStatementToBudgetCompositeEntriesResponse struct {
	BudgetCompositeEntries []budgetEntry.CompositeEntry
}
