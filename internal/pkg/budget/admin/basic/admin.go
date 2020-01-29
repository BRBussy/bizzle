package basic

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/budget"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	statementParser "github.com/BRBussy/bizzle/internal/pkg/budget/statement/parser"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
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

func (a admin) XLSXStandardBankStatementToXLSXBudgets(request *budgetAdmin.XLSXStandardBankStatementToXLSXBudgetsRequest) (*budgetAdmin.XLSXStandardBankStatementToXLSXBudgetsResponse, error) {
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
		fmt.Printf("%d %s\n%v\n", item.Year, item.Month, item.Summary)
	}

	return &budgetAdmin.XLSXStandardBankStatementToXLSXBudgetsResponse{}, nil
}

func (a admin) BudgetEntriesToBudgets(request *budgetAdmin.BudgetEntriesToBudgetsRequest) (*budgetAdmin.BudgetEntriesToBudgetResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	budgetEntryIndex := make(map[string]*budget.Budget)

	for _, entry := range request.BudgetEntries {
		budgetEntryDate, err := time.Parse(budget.DateFormat, entry.Date)
		if err != nil {
			err = bizzleException.ErrUnexpected{Reasons: []string{"could not parse budget entry date", err.Error()}}
			log.Error().Err(err)
			return nil, err
		}
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
				Summary: make(map[budgetEntry.Category]float64),
				Entries: make(map[budgetEntry.Category][]budgetEntry.Entry, 0),
			}
		}
		// update entry
		budgetEntryIndex[budgetEntryIdx].Summary[entry.Category] = budgetEntryIndex[budgetEntryIdx].Summary[entry.Category] + entry.Amount
		budgetEntryIndex[budgetEntryIdx].Entries[entry.Category] = append(budgetEntryIndex[budgetEntryIdx].Entries[entry.Category], entry)
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

func (a admin) XLSXStandardBankStatementToBudgets(request *budgetAdmin.XLSXStandardBankStatementToBudgetsRequest) (*budgetAdmin.XLSXStandardBankStatementToBudgetsResponse, error) {
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

	return &budgetAdmin.XLSXStandardBankStatementToBudgetsResponse{Budgets: budgetEntriesToBudgetsResponse.Budgets}, nil
}
