package statement

import "github.com/BRBussy/bizzle/internal/pkg/budget"

type Parser interface {
	ParseStatement(*ParseStatementRequest) (*ParseStatementResponse, error)
}

type ParseStatementRequest struct {
	Statement []byte
}

type ParseStatementResponse struct {
	Entries []budget.Entry
}
