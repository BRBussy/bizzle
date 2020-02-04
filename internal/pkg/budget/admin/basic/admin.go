package basic

import (
	budgetAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/admin"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
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

	return &budgetAdmin.GetBudgetForDateRangeResponse{}, nil
}
