package basic

import (
	budgetEntryCategoryRuleAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin"
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type admin struct {
	validator                    validationValidator.Validator
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store
}

func New(
	validator validationValidator.Validator,
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store,
) budgetEntryCategoryRuleAdmin.Admin {
	return &admin{
		validator:                    validator,
		budgetEntryCategoryRuleStore: budgetEntryCategoryRuleStore,
	}
}

func (a *admin) CreateOne(request *budgetEntryCategoryRuleAdmin.CreateOneRequest) (*budgetEntryCategoryRuleAdmin.CreateOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

}

func (a *admin) CategoriseBudgetEntry(request *budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryRequest) (*budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
}
