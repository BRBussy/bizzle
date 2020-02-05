package budget

import (
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	budgetEntryCategoryRuleAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/admin"
	budgetEntryCategoryRuleStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/rs/zerolog/log"
)

func SyncBudgetCategoryRulesForUser(
	userClaims claims.Claims,
	budgetEntryCategoryRuleAdminImp budgetEntryCategoryRuleAdmin.Admin,
	budgetEntryCategoryRuleStoreImp budgetEntryCategoryRuleStore.Store,
) error {
	// retrieve all rules owned by user
	findManyRulesResponse, err := budgetEntryCategoryRuleStoreImp.FindMany(&budgetEntryCategoryRuleStore.FindManyRequest{
		Claims: userClaims,
	})
	if err != nil {
		log.Error().Err(err).Msg("registering bizzle root user")
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
					Claims:       userClaims,
					CategoryRule: existingRule,
				}); err != nil {
					log.Error().Err(err).Msg("updating budget category rule")
					return bizzleException.ErrUnexpected{}
				}
			}
			// go to next rule
			continue nextRuleToSync
		}

		// if execution reaches here then ruleToSync does not yet exist, create it
		if _, err := budgetEntryCategoryRuleAdminImp.CreateOne(&budgetEntryCategoryRuleAdmin.CreateOneRequest{
			Claims:       userClaims,
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
		CategoryIdentifiers: []string{
			"electricity",
			"fee",
		},
		Name:   "Electricity",
		Strict: true,
	},
	{
		CategoryIdentifiers: []string{
			"wesbank",
		},
		Name: "CarRepayment",
	},
	{
		CategoryIdentifiers: []string{
			"vod",
			"prepaid",
		},
		Name:   "CellphoneAirtimeData",
		Strict: true,
	},
	{
		CategoryIdentifiers: []string{
			"telkommobi",
		},
		Name: "Internet",
	},
	{
		CategoryIdentifiers: []string{
			"disc",
			"prem",
			"medical",
		},
		Name:   "MedicalAid",
		Strict: true,
	},
	{
		CategoryIdentifiers: []string{
			"salary",
			"andile",
		},
		Name:   "Salary",
		Strict: true,
	},
}
