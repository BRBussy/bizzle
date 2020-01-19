package XLSXStandardBank

import "github.com/BRBussy/bizzle/internal/pkg/budget/statement"

type Parser struct {
}

func (p Parser) ParseStatement(*statement.ParseStatementRequest) (*statement.ParseStatementResponse, error) {
	return &statement.ParseStatementResponse{}, nil
}
