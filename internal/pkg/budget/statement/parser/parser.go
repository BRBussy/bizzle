package parser

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
)

type Parser interface {
	ParseStatementToBudgetEntries(*ParseStatementRequest) (*ParseStatementResponse, error)
}

type ParseStatementRequest struct {
	Claims    claims.Claims
	Statement []byte
}

type ParseStatementResponse struct {
	Entries []budgetEntry.Entry
}
