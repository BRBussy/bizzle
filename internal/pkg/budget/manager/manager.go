package manager

type Manager interface {
	ParseStatementXLSX(*ParseStatementXLSXRequest) (*ParseStatementXLSXResponse, error)
}

type ParseStatementXLSXRequest struct {
	XLSXStatement string
}

type ParseStatementXLSXResponse struct {
}
