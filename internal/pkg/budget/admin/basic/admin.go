package basic

import (
	"fmt"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	statementParser "github.com/BRBussy/bizzle/internal/pkg/budget/statement/parser"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type admin struct {
	validator                       validationValidator.Validator
	xlsxStandardBankStatementParser statementParser.Parser
}

func New(
	validator validationValidator.Validator,
	xlsxStandardBankStatementParser statementParser.Parser,
) budgetAdmin.Admin {
	return &admin{
		validator:                       validator,
		xlsxStandardBankStatementParser: xlsxStandardBankStatementParser,
	}
}

func (a admin) XLSXStandardBankStatementToXLSXBudget(request *budgetAdmin.XLSXStandardBankStatementToXLSXBudgetRequest) (*budgetAdmin.XLSXStandardBankStatementToXLSXBudgetResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	parseStatementResponse, err := a.xlsxStandardBankStatementParser.ParseStatement(&statementParser.ParseStatementRequest{
		Statement: request.XLSXStatement,
	})
	if err != nil {
		log.Error().Err(err).Msg("error parsing statement")
		return nil, err
	}

	for _, item := range parseStatementResponse.Entries {
		fmt.Printf("%v\n", item)
	}

	return &budgetAdmin.XLSXStandardBankStatementToXLSXBudgetResponse{}, nil
}
