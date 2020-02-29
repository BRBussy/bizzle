package basic

import (
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	budgetEntryValidator "github.com/BRBussy/bizzle/internal/pkg/budget/entry/validator"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type validator struct {
	validator                    validationValidator.Validator
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store
}

// New creates a new basic budget entry validator
func New(
	validator validationValidator.Validator,
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store,
) budgetEntryValidator.Validator {
	return &validator{
		validator:                    validator,
		budgetEntryCategoryRuleStore: budgetEntryCategoryRuleStore,
	}
}

func (v *validator) ValidateForCreate(request *budgetEntryValidator.ValidateForCreateRequest) (*budgetEntryValidator.ValidateForCreateResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetEntryValidator.ValidateForCreateResponse{}, nil
}

func (v *validator) ValidateForUpdate(request *budgetEntryValidator.ValidateForUpdateRequest) (*budgetEntryValidator.ValidateForUpdateResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetEntryValidator.ValidateForUpdateResponse{}, nil
}
