package basic

import (
	budgetConfigAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/config/admin"
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	"strings"

	budgetEntryCategoryRuleAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin"
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	validator                    validationValidator.Validator
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store
	budgetConfigAdmin            budgetConfigAdmin.Admin
}

func New(
	validator validationValidator.Validator,
	budgetEntryCategoryRuleStore budgetEntryCategoryRuleStore.Store,
	budgetConfigAdmin budgetConfigAdmin.Admin,
) budgetEntryCategoryRuleAdmin.Admin {
	return &admin{
		validator:                    validator,
		budgetEntryCategoryRuleStore: budgetEntryCategoryRuleStore,
		budgetConfigAdmin:            budgetConfigAdmin,
	}
}

func (a *admin) CreateOne(request budgetEntryCategoryRuleAdmin.CreateOneRequest) (*budgetEntryCategoryRuleAdmin.CreateOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	request.CategoryRule.ID = identifier.ID(uuid.NewV4().String())
	request.CategoryRule.OwnerID = request.Claims.ScopingID()

	if _, err := a.budgetEntryCategoryRuleStore.CreateOne(budgetEntryCategoryRuleStore.CreateOneRequest{
		CategoryRule: request.CategoryRule,
	}); err != nil {
		log.Error().Err(err).Msg("could not create budget entry category rule")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryCategoryRuleAdmin.CreateOneResponse{CategoryRule: request.CategoryRule}, nil
}

func (a *admin) UpdateOne(request budgetEntryCategoryRuleAdmin.UpdateOneRequest) (*budgetEntryCategoryRuleAdmin.UpdateOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// retrieve category rule to update
	findOneRuleResponse, err := a.budgetEntryCategoryRuleStore.FindOne(
		budgetEntryCategoryRuleStore.FindOneRequest{
			Identifier: request.CategoryRule.ID,
			Claims:     request.Claims,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve budget entry rule to update")
		return nil, bizzleException.ErrUnexpected{}
	}

	// update allowed fields
	findOneRuleResponse.CategoryRule.CategoryIdentifiers = request.CategoryRule.CategoryIdentifiers
	findOneRuleResponse.CategoryRule.Strict = request.CategoryRule.Strict
	findOneRuleResponse.CategoryRule.ExpectedAmount = request.CategoryRule.ExpectedAmount
	findOneRuleResponse.CategoryRule.ExpectedAmountPeriod = request.CategoryRule.ExpectedAmountPeriod

	// perform update
	if _, err := a.budgetEntryCategoryRuleStore.UpdateOne(
		budgetEntryCategoryRuleStore.UpdateOneRequest{
			CategoryRule: findOneRuleResponse.CategoryRule,
		},
	); err != nil {
		log.Error().Err(err).Msg("could not perform update")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryCategoryRuleAdmin.UpdateOneResponse{}, nil
}

func (a *admin) CategoriseBudgetEntry(request budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryRequest) (*budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// get users budget config to find their default other rule
	getMyConfigResponse, err := a.budgetConfigAdmin.GetMyConfig(
		budgetConfigAdmin.GetMyConfigRequest{
			Claims: request.Claims,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("unable to get budget config for user")
		return nil, bizzleException.ErrUnexpected{}
	}

	// confirm that the default other rule has been set
	if getMyConfigResponse.BudgetConfig.OtherCategoryRuleID == "" {
		return nil, budgetEntryCategoryRuleAdmin.ErrDefaultOtherBudgetEntryCategoryRuleNotSet{}
	}

	// find all category rules owned by user
	findManyBudgetCategoryRulesResponse, err := a.budgetEntryCategoryRuleStore.FindMany(
		budgetEntryCategoryRuleStore.FindManyRequest{
			Claims:   request.Claims,
			Criteria: make(criteria.Criteria, 0),
			Query:    mongo.Query{},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("could not find budget entry rules")
		return nil, bizzleException.ErrUnexpected{}
	}

	// confirm that default other rule is among those retrieved
	var otherCategoryRule budgetEntryCategoryRule.CategoryRule
	for i, bcr := range findManyBudgetCategoryRulesResponse.Records {
		if bcr.ID == getMyConfigResponse.BudgetConfig.OtherCategoryRuleID {
			otherCategoryRule = bcr
			break
		}
		if i == len(findManyBudgetCategoryRulesResponse.Records)-1 {
			// if execution reaches here then the default other rule is not found
			log.Error().Msg("other category rule not among found rules")
			return nil, bizzleException.ErrUnexpected{}
		}
	}

	// minimise and strip description
	description := strings.ToLower(strings.Trim(request.BudgetEntryDescription, " "))

	// categorisation score
	score := make(map[identifier.ID]int)

nextCategorisationRule:
	for _, rule := range findManyBudgetCategoryRulesResponse.Records {

		// if this is a strict rule score is not considered
		if rule.Strict {
			// all identifiers must be found in description
			for _, id := range rule.CategoryIdentifiers {
				if !strings.Contains(description, id) {
					// if any 1 is not found, go to next rule
					continue nextCategorisationRule
				}
			}
			// if execution reaches here then all category identifiers were found
			return &budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse{
				CategoryRule: rule,
			}, nil
		} else {
			// for non-strict rules a score is considered - i.e. count the number of identifiers found in description
			matchedIdentifiers := make([]string, 0)
			for _, id := range rule.CategoryIdentifiers {
				if strings.Contains(description, id) {
					// mark that one was found
					matchedIdentifiers = append(matchedIdentifiers, id)
				}
			}
			if len(matchedIdentifiers) > 0 {
				score[rule.ID] = len(matchedIdentifiers)
			}
		}
	}

	// if execution reaches here then we need to check the score
	highestScoreID := otherCategoryRule.ID
	highestScore := 0
	for budgetCategoryRuleID, score := range score {
		if score > highestScore {
			highestScoreID = budgetCategoryRuleID
			highestScore = score
		}
	}
	for _, rule := range findManyBudgetCategoryRulesResponse.Records {
		if rule.ID == highestScoreID {
			return &budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse{
				CategoryRule: rule,
			}, nil
		}
	}

	return &budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse{
		CategoryRule: otherCategoryRule,
	}, nil
}
