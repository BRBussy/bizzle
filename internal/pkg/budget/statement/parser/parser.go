package parser

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Parser interface {
	ParseStatementToBudgetEntries(*ParseStatementRequest) (*ParseStatementResponse, error)
}

type ParseStatementRequest struct {
	Claims    claims.Claims `validate:"required"`
	Statement []byte        `validate:"required"`
}

type ParseStatementResponse struct {
	Entries []budgetEntry.Entry
}
