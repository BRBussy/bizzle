package basic

import (
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	budgetEntryValidator "github.com/BRBussy/bizzle/internal/pkg/budget/entry/validator"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type validator struct {
	validator                    validationValidator.Validator
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store
}

// New creates a new basic budget entry validator
func New(
	validationValidator validationValidator.Validator,
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store,
) budgetEntryValidator.Validator {
	return &validator{
		validator:                    validationValidator,
		budgetEntryCategoryRuleStore: budgetEntryCategoryRuleStore,
	}
}

func (v *validator) ValidateForCreate(request *budgetEntryValidator.ValidateForCreateRequest) (*budgetEntryValidator.ValidateForCreateResponse, error) {
	if err := v.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	reasonsInvalid := make(reasonInvalid.ReasonsInvalid, 0)

	if request.BudgetEntry.CategoryRuleID != "" {
		// if category rule ID is not blank, this user should be able to retrieve it
		if _, err := v.budgetEntryCategoryRuleStore.FindOne(budgetEntryCategoryRuleStore.FindOneRequest{
			Claims:     request.Claims,
			Identifier: request.BudgetEntry.CategoryRuleID,
		}); err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				reasonsInvalid = append(
					reasonsInvalid,
					reasonInvalid.ReasonInvalid{
						Field: "categoryRuleID",
						Type:  reasonInvalid.DoesntExist,
						Help:  "must exist",
						Data:  request.BudgetEntry.CategoryRuleID,
					},
				)
			default:
				log.Error().Err(err).Msg("unable to retrieve budget category rule")
				return nil, bizzleException.ErrUnexpected{}
			}
		}
	}

	return &budgetEntryValidator.ValidateForCreateResponse{}, nil
}

func (v *validator) ValidateForUpdate(request *budgetEntryValidator.ValidateForUpdateRequest) (*budgetEntryValidator.ValidateForUpdateResponse, error) {
	if err := v.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	reasonsInvalid := make(reasonInvalid.ReasonsInvalid, 0)

	if request.BudgetEntry.CategoryRuleID != "" {
		// if category rule ID is not blank, this user should be able to retrieve it
		if _, err := v.budgetEntryCategoryRuleStore.FindOne(budgetEntryCategoryRuleStore.FindOneRequest{
			Claims:     request.Claims,
			Identifier: request.BudgetEntry.CategoryRuleID,
		}); err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				reasonsInvalid = append(
					reasonsInvalid,
					reasonInvalid.ReasonInvalid{
						Field: "categoryRuleID",
						Type:  reasonInvalid.DoesntExist,
						Help:  "must exist",
						Data:  request.BudgetEntry.CategoryRuleID,
					},
				)
			default:
				log.Error().Err(err).Msg("unable to retrieve budget category rule")
				return nil, bizzleException.ErrUnexpected{}
			}
		}
	}

	return &budgetEntryValidator.ValidateForUpdateResponse{}, nil
}
