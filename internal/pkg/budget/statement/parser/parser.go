package parser

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
)

type Parser interface {
	ParseStatement(*ParseStatementRequest) (*ParseStatementResponse, error)
}

type ParseStatementRequest struct {
	Statement []byte
}

type ParseStatementResponse struct {
	Entries []budgetEntry.Entry
}
