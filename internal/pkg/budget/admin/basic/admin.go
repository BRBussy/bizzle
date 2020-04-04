package basic

import (
	"math"

	"github.com/BRBussy/bizzle/internal/pkg/budget"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	dateTimeCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/dateTime"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type admin struct {
	validator        validationValidator.Validator
	budgetEntryStore budgetEntryStore.Store
}

// New creates a new basic budget admin
func New(
	validator validationValidator.Validator,
	budgetEntryStore budgetEntryStore.Store,
) budgetAdmin.Admin {
	return &admin{
		budgetEntryStore: budgetEntryStore,
		validator:        validator,
	}
}

func (a *admin) GetBudgetForMonthInYear(request budgetAdmin.GetBudgetForMonthInYearRequest) (*budgetAdmin.GetBudgetForMonthInYearResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetAdmin.GetBudgetForMonthInYearResponse{}, nil
}

func (a *admin) GetBudgetForDateRange(request budgetAdmin.GetBudgetForDateRangeRequest) (*budgetAdmin.GetBudgetForDateRangeResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// retrieve all budget entries in date range
	findManyBudgetEntriesResponse, err := a.budgetEntryStore.FindManyComposite(
		budgetEntryStore.FindManyCompositeRequest{
			Claims: request.Claims,
			Criteria: criteria.Criteria{
				dateTimeCriterion.Range{
					Field: "date",
					Start: dateTimeCriterion.RangeValue{
						Date:      request.StartDate,
						Inclusive: true,
						Ignore:    false,
					},
					End: dateTimeCriterion.RangeValue{
						Date:      request.EndDate,
						Inclusive: false,
					},
				},
			},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve budget entries")
		return nil, bizzleException.ErrUnexpected{Reasons: []string{"could not retrieve budget entries"}}
	}

	// create and populate budget
	newBudget := budget.Budget{
		StartDate: request.StartDate,
		EndDate:   request.EndDate,
		Summary:   make(map[string]float64),
		Entries:   make(map[string][]budgetEntry.Entry),
	}
	for _, be := range findManyBudgetEntriesResponse.Records {
		// sum amounts of all the entries with the same category rule
		newBudget.Summary[be.CategoryRule.Name] = newBudget.Summary[be.CategoryRule.Name] + be.Amount

		// put together all entries with the same category rule
		newBudget.Entries[be.CategoryRule.Name] = append(newBudget.Entries[be.CategoryRule.Name], be.Entry)
	}

	// perform rounding on summary
	for summaryKey, value := range newBudget.Summary {
		if value > 0 {
			newBudget.TotalIn.Actual += value
		} else {
			newBudget.TotalOut.Actual -= value
		}
		newBudget.Summary[summaryKey] = math.Round(value*100) / 100
	}
	newBudget.Net = newBudget.TotalIn.Actual - newBudget.TotalOut.Actual
	newBudget.Net = math.Round(newBudget.Net*100) / 100
	newBudget.TotalIn.Actual = math.Round(newBudget.TotalIn.Actual*100) / 100
	newBudget.TotalOut.Actual = math.Round(newBudget.TotalOut.Actual*100) / 100

	return &budgetAdmin.GetBudgetForDateRangeResponse{Budget: newBudget}, nil
}
