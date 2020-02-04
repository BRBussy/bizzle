package basic

import (
	"fmt"
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
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

func New(
	validator validationValidator.Validator,
	budgetEntryStore budgetEntryStore.Store,
) budgetAdmin.Admin {
	return &admin{
		budgetEntryStore: budgetEntryStore,
		validator:        validator,
	}
}

func (a *admin) GetBudgetForMonthInYear(request *budgetAdmin.GetBudgetForMonthInYearRequest) (*budgetAdmin.GetBudgetForMonthInYearResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetAdmin.GetBudgetForMonthInYearResponse{}, nil
}

func (a *admin) GetBudgetForDateRange(request *budgetAdmin.GetBudgetForDateRangeRequest) (*budgetAdmin.GetBudgetForDateRangeResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// retrieve all budget entries in date range
	findManyBudgetEntriesResponse, err := a.budgetEntryStore.FindMany(&budgetEntryStore.FindManyRequest{
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
	})
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve budget entries")
		return nil, bizzleException.ErrUnexpected{Reasons: []string{"could not retrieve budget entries"}}
	}

	fmt.Println(findManyBudgetEntriesResponse.Total)

	return &budgetAdmin.GetBudgetForDateRangeResponse{}, nil
}
