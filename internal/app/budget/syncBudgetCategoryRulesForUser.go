package budget

import (
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	budgetEntryCategoryRuleAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin"
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

func SyncBudgetCategoryRulesForUser(
	userID identifier.ID,
	budgetEntryCategoryRuleAdminImp budgetEntryCategoryRuleAdmin.Admin,
	budgetEntryCategoryRuleStoreImp budgetEntryCategoryRuleStore.Store,
	userStoreImp userStore.Store,
) error {
	// retrieve the user
	findOneUserResponse, err := userStoreImp.FindOne(&userStore.FindOneRequest{
		Identifier: userID,
	})
	if err != nil {
		log.Error().Err(err).Msg("could not retrieve user")
		return err
	}

	// retrieve all rules owned by user
	findManyRulesResponse, err := budgetEntryCategoryRuleStoreImp.FindMany(&budgetEntryCategoryRuleStore.FindManyRequest{
		Criteria: make(criteria.Criteria, 0),
	})
	if err != nil {
		log.Error().Err(err).Msg("retrieving all budget entry category rules owned by user")
		return bizzleException.ErrUnexpected{}
	}

	// for every rule to sync
nextRuleToSync:
	for _, ruleToSync := range categoryRulesToSyncForUser {
		// look to see if the rule already exists
		for _, existingRule := range findManyRulesResponse.Records {
			// if it already exists (has the same name) AND
			// it has changed, update it
			if ruleToSync.Name == existingRule.Name &&
				!budgetEntryCategoryRule.CompareCategoryRules(ruleToSync, existingRule) {
				existingRule.Strict = ruleToSync.Strict
				existingRule.CategoryIdentifiers = ruleToSync.CategoryIdentifiers
				if _, err := budgetEntryCategoryRuleAdminImp.UpdateOne(&budgetEntryCategoryRuleAdmin.UpdateOneRequest{
					Claims:       claims.Login{UserID: userID},
					CategoryRule: existingRule,
				}); err != nil {
					log.Error().Err(err).Msg("updating budget category rule")
					return bizzleException.ErrUnexpected{}
				}
				// go to next rule
				continue nextRuleToSync
			}
		}

		// if execution reaches here then ruleToSync does not yet exist, create it
		ruleToSync.OwnerID = findOneUserResponse.User.ID
		if _, err := budgetEntryCategoryRuleAdminImp.CreateOne(&budgetEntryCategoryRuleAdmin.CreateOneRequest{
			Claims:       claims.Login{UserID: userID},
			CategoryRule: ruleToSync,
		}); err != nil {
			log.Error().Err(err).Msg("creating budget category rule")
			return bizzleException.ErrUnexpected{}
		}
	}

	return nil
}

var categoryRulesToSyncForUser = []budgetEntryCategoryRule.CategoryRule{
	{
		Name: "Electricity",
		CategoryIdentifiers: []string{
			"electricity",
			"fee",
		},
		Strict: true,
	},
	{
		Name: "CarRepayment",
		CategoryIdentifiers: []string{
			"wesbank",
		},
	},
	{
		Name: "CellphoneAirtimeData",
		CategoryIdentifiers: []string{
			"vod",
			"prepaid",
		},
		Strict: true,
	},
	{
		Name: "Internet",
		CategoryIdentifiers: []string{
			"telkommobi",
		},
	},
	{
		Name: "MedicalAid",
		CategoryIdentifiers: []string{
			"disc",
			"prem",
			"medical",
		},
		Strict: true,
	},
	{
		Name: "Salary",
		CategoryIdentifiers: []string{
			"salary",
			"andile",
		},
		Strict: true,
	},
}
