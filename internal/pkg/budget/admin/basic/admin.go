package basic

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	statementParser "github.com/BRBussy/bizzle/internal/pkg/budget/statement/parser"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	"time"
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

	// parse standard bank statement
	parseStatementResponse, err := a.xlsxStandardBankStatementParser.ParseStatement(&statementParser.ParseStatementRequest{
		Statement: request.XLSXStatement,
	})
	if err != nil {
		log.Error().Err(err).Msg("error parsing statement")
		return nil, err
	}

	// process resultant budget entries into budgets
	budgetEntriesToBudgetsResponse, err := a.BudgetEntriesToBudgets(&budgetAdmin.BudgetEntriesToBudgetsRequest{
		BudgetEntries: parseStatementResponse.Entries,
	})
	if err != nil {
		log.Error().Err(err).Msg("error processing budget entries into budgets")
		return nil, err
	}

	for _, item := range budgetEntriesToBudgetsResponse.Budgets {
		fmt.Printf("%v\n", item)
	}

	return &budgetAdmin.XLSXStandardBankStatementToXLSXBudgetResponse{}, nil
}

func (a admin) BudgetEntriesToBudgets(request *budgetAdmin.BudgetEntriesToBudgetsRequest) (*budgetAdmin.BudgetEntriesToBudgetResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	budgetEntryIndex := make(map[string]*budget.Budget)

	for _, budgetEntry := range request.BudgetEntries {
		budgetEntryDate := time.Unix(budgetEntry.Date, 0)
		budgetEntryIdx := fmt.Sprintf(
			"%d-%s",
			budgetEntryDate.Year(),
			budgetEntryDate.Month(),
		)
		// if an entry has not yet been made in index, make one
		if _, found := budgetEntryIndex[budgetEntryIdx]; !found {
			budgetEntryIndex[budgetEntryIdx] = &budget.Budget{
				Month:   budgetEntryDate.Month().String(),
				Year:    budgetEntryDate.Year(),
				Summary: make(map[budget.Category]float64),
				Entries: make(map[budget.Category][]budget.Entry, 0),
			}
		}
		// update entry
	}

	budgets := make([]budget.Budget, 0)
	for _, budgetPtr := range budgetEntryIndex {
		budgets = append(budgets, *budgetPtr)
	}
	return &budgetAdmin.BudgetEntriesToBudgetResponse{Budgets: budgets}, nil
}

func (a admin) BudgetToXLSX(request *budgetAdmin.BudgetToXLSXRequest) (*budgetAdmin.BudgetToXLSXResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetAdmin.BudgetToXLSXResponse{}, nil
}
