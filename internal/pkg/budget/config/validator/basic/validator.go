package basic

import (
	budgetConfigValidator "github.com/BRBussy/bizzle/internal/pkg/budget/config/validator"
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
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

// New creates a new basic budget config validator
func New(
	validationValidator validationValidator.Validator,
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store,
) budgetConfigValidator.Validator {
	return &validator{
		validator:                    validationValidator,
		budgetEntryCategoryRuleStore: budgetEntryCategoryRuleStore,
	}
}

func (v *validator) ValidateForCreate(request budgetConfigValidator.ValidateForCreateRequest) (*budgetConfigValidator.ValidateForCreateResponse, error) {
	if err := v.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	reasonsInvalid := make(reasonInvalid.ReasonsInvalid, 0)

	if request.BudgetConfig.OtherCategoryRuleID != "" {
		// if other category rule ID is not blank, this user should be able to retrieve it
		if _, err := v.budgetEntryCategoryRuleStore.FindOne(
			budgetEntryCategoryRuleStore.FindOneRequest{
				Claims:     request.Claims,
				Identifier: request.BudgetConfig.OtherCategoryRuleID,
			},
		); err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				reasonsInvalid = append(
					reasonsInvalid,
					reasonInvalid.ReasonInvalid{
						Field: "otherCategoryRuleID",
						Type:  reasonInvalid.DoesntExist,
						Help:  "must exist",
						Data:  request.BudgetConfig.OtherCategoryRuleID,
					},
				)
			default:
				log.Error().Err(err).Msg("unable to retrieve budget config other category rule")
				return nil, bizzleException.ErrUnexpected{}
			}
		}
	}

	if request.BudgetConfig.SummaryDatePeriodCategoryRuleID != "" {
		// if summary date period category rule ID is not blank, this user should be able to retrieve it
		if _, err := v.budgetEntryCategoryRuleStore.FindOne(
			budgetEntryCategoryRuleStore.FindOneRequest{
				Claims:     request.Claims,
				Identifier: request.BudgetConfig.SummaryDatePeriodCategoryRuleID,
			},
		); err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				reasonsInvalid = append(
					reasonsInvalid,
					reasonInvalid.ReasonInvalid{
						Field: "summaryDatePeriodCategoryRuleID",
						Type:  reasonInvalid.DoesntExist,
						Help:  "must exist",
						Data:  request.BudgetConfig.SummaryDatePeriodCategoryRuleID,
					},
				)
			default:
				log.Error().Err(err).Msg("unable to retrieve budget config summary date period category rule")
				return nil, bizzleException.ErrUnexpected{}
			}
		}
	}

	return &budgetConfigValidator.ValidateForCreateResponse{}, nil
}

func (v *validator) ValidateForUpdate(request budgetConfigValidator.ValidateForUpdateRequest) (*budgetConfigValidator.ValidateForUpdateResponse, error) {
	if err := v.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	reasonsInvalid := make(reasonInvalid.ReasonsInvalid, 0)

	if request.BudgetConfig.OtherCategoryRuleID != "" {
		// if category rule ID is not blank, this user should be able to retrieve it
		if _, err := v.budgetEntryCategoryRuleStore.FindOne(
			budgetEntryCategoryRuleStore.FindOneRequest{
				Claims:     request.Claims,
				Identifier: request.BudgetConfig.OtherCategoryRuleID,
			},
		); err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				reasonsInvalid = append(
					reasonsInvalid,
					reasonInvalid.ReasonInvalid{
						Field: "otherCategoryRuleID",
						Type:  reasonInvalid.DoesntExist,
						Help:  "must exist",
						Data:  request.BudgetConfig.OtherCategoryRuleID,
					},
				)
			default:
				log.Error().Err(err).Msg("unable to retrieve budget category rule")
				return nil, bizzleException.ErrUnexpected{}
			}
		}
	}

	if request.BudgetConfig.SummaryDatePeriodCategoryRuleID != "" {
		// if summary date period category rule ID is not blank, this user should be able to retrieve it
		if _, err := v.budgetEntryCategoryRuleStore.FindOne(
			budgetEntryCategoryRuleStore.FindOneRequest{
				Claims:     request.Claims,
				Identifier: request.BudgetConfig.SummaryDatePeriodCategoryRuleID,
			},
		); err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				reasonsInvalid = append(
					reasonsInvalid,
					reasonInvalid.ReasonInvalid{
						Field: "summaryDatePeriodCategoryRuleID",
						Type:  reasonInvalid.DoesntExist,
						Help:  "must exist",
						Data:  request.BudgetConfig.SummaryDatePeriodCategoryRuleID,
					},
				)
			default:
				log.Error().Err(err).Msg("unable to retrieve budget config summary date period category rule")
				return nil, bizzleException.ErrUnexpected{}
			}
		}
	}

	return &budgetConfigValidator.ValidateForUpdateResponse{}, nil
}
