package basic

import (
	"math"

	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	statementParser "github.com/BRBussy/bizzle/internal/pkg/budget/statement/parser"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	validator                       validationValidator.Validator
	budgetEntryStore                budgetEntryStore.Store
	xlsxStandardBankStatementParser statementParser.Parser
}

// New creates a new basic budget entry admin
func New(
	validator validationValidator.Validator,
	budgetEntryStore budgetEntryStore.Store,
	xlsxStandardBankStatementParser statementParser.Parser,
) budgetEntryAdmin.Admin {
	return &admin{
		budgetEntryStore:                budgetEntryStore,
		validator:                       validator,
		xlsxStandardBankStatementParser: xlsxStandardBankStatementParser,
	}
}

func (a *admin) CreateMany(request *budgetEntryAdmin.CreateManyRequest) (*budgetEntryAdmin.CreateManyResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	for entryIdx := range request.BudgetEntries {
		// set ID and onwerID
		request.BudgetEntries[entryIdx].OwnerID = request.Claims.ScopingID()
		request.BudgetEntries[entryIdx].ID = identifier.ID(uuid.NewV4().String())

		// round off to 2 units
		request.BudgetEntries[entryIdx].Amount = math.Round(request.BudgetEntries[entryIdx].Amount*100) / 100
	}

	if _, err := a.budgetEntryStore.CreateMany(&budgetEntryStore.CreateManyRequest{
		Entries: request.BudgetEntries,
	}); err != nil {
		log.Error().Err(err).Msg("could not create many budget entries")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryAdmin.CreateManyResponse{}, nil
}

func (a *admin) DuplicateCheck(request *budgetEntryAdmin.DuplicateCheckRequest) (*budgetEntryAdmin.DuplicateCheckResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetEntryAdmin.DuplicateCheckResponse{
		Uniques:             request.BudgetEntries,
		ExactDuplicates:     make([]budgetEntry.Entry, 0),
		SuspectedDuplicates: make([]budgetEntry.Entry, 0),
	}, nil
}

func (a *admin) XLSXStandardBankStatementToBudgetEntries(
	request *budgetEntryAdmin.XLSXStandardBankStatementToBudgetEntriesRequest,
) (*budgetEntryAdmin.XLSXStandardBankStatementToBudgetEntriesResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// parse standard bank statement
	parseStatementToBudgetEntriesResponse, err := a.xlsxStandardBankStatementParser.ParseStatementToBudgetEntries(
		&statementParser.ParseStatementToBudgetEntriesRequest{
			Claims:    request.Claims,
			Statement: request.XLSXStatement,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("error parsing statement to budget entries")
		return nil, err
	}

	return &budgetEntryAdmin.XLSXStandardBankStatementToBudgetEntriesResponse{
		BudgetEntries: parseStatementToBudgetEntriesResponse.BudgetEntries,
	}, nil
}
