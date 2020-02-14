package parser

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Parser interface {
	ParseStatementToBudgetEntries(*ParseStatementToBudgetEntriesRequest) (*ParseStatementToBudgetEntriesResponse, error)
}

type ParseStatementToBudgetEntriesRequest struct {
	Claims    claims.Claims `validate:"required"`
	Statement []byte        `validate:"required"`
}

type ParseStatementToBudgetEntriesResponse struct {
	BudgetEntries []budgetEntry.Entry
}
