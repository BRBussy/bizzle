package basic

import (
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	budgetEntryCategoryRuleAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin"
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"strings"
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

	request.CategoryRule.ID = identifier.ID(uuid.NewV4().String())
	request.CategoryRule.OwnerID = request.Claims.ScopingID()

	if _, err := a.budgetEntryCategoryRuleStore.CreateOne(&budgetEntryCategoryRuleStore.CreateOneRequest{
		CategoryRule: request.CategoryRule,
	}); err != nil {
		log.Error().Err(err).Msg("could not create budget entry category rule")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &budgetEntryCategoryRuleAdmin.CreateOneResponse{CategoryRule: request.CategoryRule}, nil
}

func (a *admin) UpdateOne(request *budgetEntryCategoryRuleAdmin.UpdateOneRequest) (*budgetEntryCategoryRuleAdmin.UpdateOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	findOneRuleResponse, err := a.budgetEntryCategoryRuleStore.FindOne(&budgetEntryCategoryRuleStore.FindOneRequest{
		Claims:     request.Claims,
		Identifier: request.CategoryRule.ID,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve budget entry rule to update")
		return nil, bizzleException.ErrUnexpected{}
	}

	findOneRuleResponse.CategoryRule.CategoryIdentifiers = request.CategoryRule.CategoryIdentifiers
	findOneRuleResponse.CategoryRule.Strict = request.CategoryRule.Strict

	return &budgetEntryCategoryRuleAdmin.UpdateOneResponse{}, nil
}

func (a *admin) CategoriseBudgetEntry(request *budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryRequest) (*budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// find all category rules owned by user
	findManyResponse, err := a.budgetEntryCategoryRuleStore.FindMany(&budgetEntryCategoryRuleStore.FindManyRequest{
		Claims:   request.Claims,
		Criteria: make(criteria.Criteria, 0),
		Query:    mongo.Query{},
	})
	if err != nil {
		log.Error().Err(err).Msg("could not find budget entry rules")
		return nil, bizzleException.ErrUnexpected{}
	}

	// minimise and strip description
	description := strings.ToLower(strings.Trim(request.BudgetEntryDescription, " "))

nextCategorisationRule:
	for _, rule := range findManyResponse.Records {
		if rule.Strict {
			// all identifiers must be found in description
			for _, id := range rule.CategoryIdentifiers {
				if !strings.Contains(description, id.String()) {
					// if any 1 is not found, go to next rule
					continue nextCategorisationRule
				}
			}
			// if execution reaches here then all category identifiers were found
			return &budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse{CategoryID: rule.ID}, nil
		} else {
			// any identifiers can be found in description
			matchedIdentifiers := make([]budgetEntryCategoryRule.CategoryIdentifier, 0)
			for _, id := range rule.CategoryIdentifiers {
				if strings.Contains(description, id.String()) {
					// mark that one was found
					matchedIdentifiers = append(matchedIdentifiers, id)
				}
			}
			if len(matchedIdentifiers) > 0 {
				return &budgetEntryCategoryRuleAdmin.CategoriseBudgetEntryResponse{CategoryID: rule.ID}, nil
			}
		}
	}

	return nil, budgetEntryCategoryRule.ErrCouldNotClassify{Reason: "not match"}
}
