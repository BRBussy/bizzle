package basic

import (
	"fmt"
	budgetEntryIgnored "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored"
	budgetEntryIgnoredAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored/admin"
	budgetEntryIgnoredStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"math"
	"strings"
	"time"

	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	budgetEntryValidator "github.com/BRBussy/bizzle/internal/pkg/budget/entry/validator"
	statementParser "github.com/BRBussy/bizzle/internal/pkg/budget/statement/parser"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	dateTimeCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/dateTime"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	validator                       validationValidator.Validator
	budgetEntryStore                budgetEntryStore.Store
	budgetEntryValidator            budgetEntryValidator.Validator
	xlsxStandardBankStatementParser statementParser.Parser
	budgetEntryIgnoredAdmin         budgetEntryIgnoredAdmin.Admin
	budgetEntryIgnoredStore         budgetEntryIgnoredStore.Store
}

// New creates a new basic budget entry admin
func New(
	validator validationValidator.Validator,
	budgetEntryStore budgetEntryStore.Store,
	budgetEntryValidator budgetEntryValidator.Validator,
	xlsxStandardBankStatementParser statementParser.Parser,
	budgetEntryIgnoredAdmin budgetEntryIgnoredAdmin.Admin,
	budgetEntryIgnoredStore budgetEntryIgnoredStore.Store,
) budgetEntryAdmin.Admin {
	return &admin{
		budgetEntryStore:                budgetEntryStore,
		validator:                       validator,
		budgetEntryValidator:            budgetEntryValidator,
		xlsxStandardBankStatementParser: xlsxStandardBankStatementParser,
		budgetEntryIgnoredStore:         budgetEntryIgnoredStore,
		budgetEntryIgnoredAdmin:         budgetEntryIgnoredAdmin,
	}
}

func (a *admin) CreateOne(request budgetEntryAdmin.CreateOneRequest) (*budgetEntryAdmin.CreateOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	request.BudgetEntry.OwnerID = request.Claims.ScopingID()
	request.BudgetEntry.ID = identifier.ID(uuid.NewV4().String())

	// round off to 2 units
	request.BudgetEntry.Amount = math.Round(request.BudgetEntry.Amount*100) / 100

	// validate the entry for create
	validateForCreateResponse, err := a.budgetEntryValidator.ValidateForCreate(&budgetEntryValidator.ValidateForCreateRequest{
		Claims:      request.Claims,
		BudgetEntry: request.BudgetEntry,
	})
	if err != nil {
		log.Error().Err(err).Msg("error validating entry for create")
		return nil, bizzleException.ErrUnexpected{}
	}
	if len(validateForCreateResponse.ReasonsInvalid) > 0 {
		return nil, budgetEntry.ErrInvalidEntry{ReasonsInvalid: validateForCreateResponse.ReasonsInvalid}
	}

	// perform creation
	if _, err := a.budgetEntryStore.CreateOne(budgetEntryStore.CreateOneRequest{
		Entry: request.BudgetEntry,
	}); err != nil {
		log.Error().Err(err).Msg("error creating budget entry")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryAdmin.CreateOneResponse{BudgetEntry: request.BudgetEntry}, nil
}

func (a *admin) CreateMany(request budgetEntryAdmin.CreateManyRequest) (*budgetEntryAdmin.CreateManyResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	for entryIdx := range request.BudgetEntries {
		// set ID and ownerID
		request.BudgetEntries[entryIdx].OwnerID = request.Claims.ScopingID()
		request.BudgetEntries[entryIdx].ID = identifier.ID(uuid.NewV4().String())

		// round off to 2 units
		request.BudgetEntries[entryIdx].Amount = math.Round(request.BudgetEntries[entryIdx].Amount*100) / 100

		// validate the entry for create
		validateForUpdateResponse, err := a.budgetEntryValidator.ValidateForCreate(&budgetEntryValidator.ValidateForCreateRequest{
			Claims:      request.Claims,
			BudgetEntry: request.BudgetEntries[entryIdx],
		})
		if err != nil {
			log.Error().Err(err).Msg("error validating entry for create")
			return nil, bizzleException.ErrUnexpected{}
		}

		// check if there are any reasons that the entry is invalid
		if len(validateForUpdateResponse.ReasonsInvalid) > 0 {
			return nil, budgetEntry.ErrInvalidEntry{
				ReasonsInvalid: validateForUpdateResponse.ReasonsInvalid,
			}
		}
	}

	if _, err := a.budgetEntryStore.CreateMany(budgetEntryStore.CreateManyRequest{
		Entries: request.BudgetEntries,
	}); err != nil {
		log.Error().Err(err).Msg("could not create many budget entries")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryAdmin.CreateManyResponse{}, nil
}

func (a *admin) DuplicateCheck(request budgetEntryAdmin.DuplicateCheckRequest) (*budgetEntryAdmin.DuplicateCheckResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// find range of category entries being checked
	var earliestDate time.Time = time.Now()
	var latestDate time.Time = time.Now()
	// for every that needs to be part of the duplicate check ...
	for _, entryToCheck := range request.BudgetEntries {
		// if it is before the listed earliest date
		if earliestDate.After(entryToCheck.Date) {
			// update earliest date to this entry's date
			earliestDate = entryToCheck.Date
		}
		// if it is after the latest date
		if latestDate.Before(entryToCheck.Date) {
			// update latest date to this entry's date
			latestDate = entryToCheck.Date
		}
	}

	// search for all budget entries in this date range
	findManyResponse, err := a.budgetEntryStore.FindMany(budgetEntryStore.FindManyRequest{
		Criteria: criteria.Criteria{
			dateTimeCriterion.Range{
				Field: "date",
				Start: dateTimeCriterion.RangeValue{
					Date:      earliestDate,
					Inclusive: true,
				},
				End: dateTimeCriterion.RangeValue{
					Date:      latestDate,
					Inclusive: true,
				},
			},
		},
		Claims: request.Claims,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not search for budget entries for duplicate check")
		return nil, bizzleException.ErrUnexpected{}
	}

	// exact duplicates are two entries that match exactly
	exactDuplicates := make([]budgetEntryAdmin.DuplicateEntries, 0)
	// suspected duplicates are two entries that are on the same date with same amount
	suspectedDuplicates := make([]budgetEntryAdmin.DuplicateEntries, 0)
	uniques := make([]budgetEntry.Entry, 0)

	// for every new entry to import...
nextEntryToImport:
	for _, entryToImport := range request.BudgetEntries {

		// check to see if it is an exact suspected duplicate of any of the existing items
		for entryIdx, existingEntry := range findManyResponse.Records {
			if existingEntry.ExactDuplicate(entryToImport) {
				// if one is found, add it to list
				exactDuplicates = append(
					exactDuplicates,
					budgetEntryAdmin.DuplicateEntries{
						Existing: existingEntry,
						New:      entryToImport,
					},
				)
				// and remove it from entry records, do not let it take part in future checks
				removeBudgetEntry(&findManyResponse.Records, entryIdx)
				continue nextEntryToImport
			}
		}

		// check to see if it a suspected duplicated of any of the existing items
		for _, existingEntry := range findManyResponse.Records {
			if existingEntry.SuspectedDuplicate(entryToImport) {
				suspectedDuplicates = append(
					suspectedDuplicates,
					budgetEntryAdmin.DuplicateEntries{
						Existing: existingEntry,
						New:      entryToImport,
					},
				)
				continue nextEntryToImport
			}
		}

		// if execution reaches here, the entry is neither an exact nor suspected duplicate
		// assume it is unique
		uniques = append(uniques, entryToImport)
	}

	return &budgetEntryAdmin.DuplicateCheckResponse{
		Uniques:             uniques,
		ExactDuplicates:     exactDuplicates,
		SuspectedDuplicates: suspectedDuplicates,
	}, nil
}

func (a *admin) XLSXStandardBankStatementToBudgetEntries(
	request budgetEntryAdmin.XLSXStandardBankStatementToBudgetEntriesRequest,
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

func removeBudgetEntry(budgetEntries *[]budgetEntry.Entry, idxToRemove int) {
	(*budgetEntries)[len(*budgetEntries)-1], (*budgetEntries)[idxToRemove] = (*budgetEntries)[idxToRemove], (*budgetEntries)[len(*budgetEntries)-1]
	*budgetEntries = (*budgetEntries)[:len(*budgetEntries)-1]
}

func (a *admin) UpdateOne(request budgetEntryAdmin.UpdateOneRequest) (*budgetEntryAdmin.UpdateOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// validate the entry for update
	validateForUpdateResponse, err := a.budgetEntryValidator.ValidateForUpdate(&budgetEntryValidator.ValidateForUpdateRequest{
		Claims:      request.Claims,
		BudgetEntry: request.BudgetEntry,
	})
	if err != nil {
		log.Error().Err(err).Msg("error validating entry for update")
		return nil, bizzleException.ErrUnexpected{}
	}

	// check if there are any reasons that the entry is invalid
	if len(validateForUpdateResponse.ReasonsInvalid) > 0 {
		return nil, budgetEntry.ErrInvalidEntry{
			ReasonsInvalid: validateForUpdateResponse.ReasonsInvalid,
		}
	}

	// retrieve the entry that needs to be updated
	findOneEntryResponse, err := a.budgetEntryStore.FindOne(budgetEntryStore.FindOneRequest{
		Claims:     request.Claims,
		Identifier: request.BudgetEntry.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve budget entry to be updated")
		return nil, bizzleException.ErrUnexpected{}
	}

	// set fields that are not allowed to be updated
	request.BudgetEntry.OwnerID = findOneEntryResponse.Entry.OwnerID

	// perform update
	if _, err := a.budgetEntryStore.UpdateOne(budgetEntryStore.UpdateOneRequest{
		Claims: request.Claims,
		Entry:  request.BudgetEntry,
	}); err != nil {
		log.Error().Err(err).Msg("could not update budget entry")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryAdmin.UpdateOneResponse{}, nil
}

func (a *admin) UpdateMany(request budgetEntryAdmin.UpdateManyRequest) (*budgetEntryAdmin.UpdateManyResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// validate and confirm retrieval of every entry being updated
	for entryToUpdateIdx := range request.BudgetEntries {
		// validate the entry for update
		validateForUpdateResponse, err := a.budgetEntryValidator.ValidateForUpdate(&budgetEntryValidator.ValidateForUpdateRequest{
			Claims:      request.Claims,
			BudgetEntry: request.BudgetEntries[entryToUpdateIdx],
		})
		if err != nil {
			log.Error().Err(err).Msg("error validating entry for update")
			return nil, bizzleException.ErrUnexpected{}
		}

		// check if there are any reasons that the entry is invalid
		if len(validateForUpdateResponse.ReasonsInvalid) > 0 {
			return nil, budgetEntry.ErrInvalidEntry{
				ReasonsInvalid: validateForUpdateResponse.ReasonsInvalid,
			}
		}

		// retrieve the entry that needs to be updated
		findOneEntryResponse, err := a.budgetEntryStore.FindOne(budgetEntryStore.FindOneRequest{
			Claims:     request.Claims,
			Identifier: request.BudgetEntries[entryToUpdateIdx].ID,
		})
		if err != nil {
			log.Error().Err(err).Msg("could not retrieve budget entry to be updated")
			return nil, bizzleException.ErrUnexpected{}
		}

		// set fields that are not allowed to be updated
		request.BudgetEntries[entryToUpdateIdx].OwnerID = findOneEntryResponse.Entry.OwnerID
	}

	// perform updates
	for _, entryToUpdate := range request.BudgetEntries {
		// perform update
		if _, err := a.budgetEntryStore.UpdateOne(budgetEntryStore.UpdateOneRequest{
			Claims: request.Claims,
			Entry:  entryToUpdate,
		}); err != nil {
			log.Error().Err(err).Msg("could not update budget entry")
			return nil, bizzleException.ErrUnexpected{}
		}
	}

	return &budgetEntryAdmin.UpdateManyResponse{}, nil
}

func (a *admin) DeleteOne(request budgetEntryAdmin.DeleteOneRequest) (*budgetEntryAdmin.DeleteOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if _, err := a.budgetEntryStore.DeleteOne(
		budgetEntryStore.DeleteOneRequest{
			Claims:     request.Claims,
			Identifier: request.Identifier,
		},
	); err != nil {
		log.Error().Err(err).Msg("could not delete budget entry")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryAdmin.DeleteOneResponse{}, nil
}

func (a *admin) IgnoreOne(request budgetEntryAdmin.IgnoreOneRequest) (*budgetEntryAdmin.IgnoreOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if _, err := a.budgetEntryIgnoredAdmin.CreateOne(
		budgetEntryIgnoredAdmin.CreateOneRequest{
			Claims: request.Claims,
			Ignored: budgetEntryIgnored.Ignored{
				Description: fmt.Sprintf(
					"%s-%s",
					request.BudgetEntry.Date.Format("010206"),
					strings.ToLower(request.BudgetEntry.Description),
				),
			},
		},
	); err != nil {
		log.Error().Err(err).Msg("unable to create ignored")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryAdmin.IgnoreOneResponse{}, nil
}

func (a *admin) RecogniseOne(request budgetEntryAdmin.RecogniseOneRequest) (*budgetEntryAdmin.RecogniseOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetEntryAdmin.RecogniseOneResponse{}, nil
}

func (a *admin) IgnoredCheck(request budgetEntryAdmin.IgnoredCheckRequest) (*budgetEntryAdmin.IgnoredCheckResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	findManyIgnored, err := a.budgetEntryIgnoredStore.FindMany(
		budgetEntryIgnoredStore.FindManyRequest{
			Claims:   request.Claims,
			Criteria: make(criteria.Criteria, 0),
			Query:    mongo.Query{},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to find many ignored")
		return nil, bizzleException.ErrUnexpected{}
	}

	ignoredBudgetEntries := make([]budgetEntryIgnored.Ignored, 0)
NextEntry:
	for _, entry := range request.BudgetEntries {
		for _, ignored := range findManyIgnored.Records {
			if ignored.Description == fmt.Sprintf(
				"%s-%s",
				entry.Date.Format("010206"),
				strings.ToLower(entry.Description),
			) {
				ignoredBudgetEntries = append(
					ignoredBudgetEntries,
					ignored,
				)
				continue NextEntry
			}
		}
	}

	return &budgetEntryAdmin.IgnoredCheckResponse{
		Ignored: ignoredBudgetEntries,
	}, nil
}
