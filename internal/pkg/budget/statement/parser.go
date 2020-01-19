package statement

type Parser interface {
	ParseStatement(*ParseStatementRequest) (*ParseStatementResponse, error)
}

type ParseStatementRequest struct {
	XLSXStatement string
}

type ParseStatementResponse struct {
}
