package statement

type Parser interface {
	ParseStatement(*ParseStatementRequest) (*ParseStatementResponse, error)
}

type ParseStatementRequest struct {
	Statement []byte
}

type ParseStatementResponse struct {
}
