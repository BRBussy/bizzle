package budget

import (
	"fmt"
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
	log.Info().Msg("Running SyncBudgetCategoryRulesForUser")
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
	log.Info().Msg(fmt.Sprintf("User has %d budget category rules", len(findManyRulesResponse.Records)))

	// for every rule to sync
nextRuleToSync:
	for ruleToSyncIdx, ruleToSync := range categoryRulesToSyncForUser {
		log.Info().Msg(fmt.Sprintf(
			"--- [%d/%d] Check Rule to sync %s ---",
			ruleToSyncIdx+1,
			len(categoryRulesToSyncForUser),
			ruleToSync.Name,
		))
		// look to see if the rule already exists
		for _, existingRule := range findManyRulesResponse.Records {

			// check if the rule already exists (has the same name)
			if ruleToSync.Name == existingRule.Name {
				log.Info().Msg(fmt.Sprintf("  Rule %s already exists", ruleToSync.Name))

				// set ids to prevent false positive comparisons
				ruleToSync.ID = existingRule.ID
				ruleToSync.OwnerID = existingRule.OwnerID

				// if it does exists, check if the rule has changed
				if budgetEntryCategoryRule.CompareCategoryRules(ruleToSync, existingRule) {
					log.Info().Msg("  no change")
				} else {
					log.Info().Msg(fmt.Sprintf("  Rule %s has changed - Updating", ruleToSync.Name))
					existingRule.Strict = ruleToSync.Strict
					existingRule.CategoryIdentifiers = ruleToSync.CategoryIdentifiers
					if _, err := budgetEntryCategoryRuleAdminImp.UpdateOne(&budgetEntryCategoryRuleAdmin.UpdateOneRequest{
						Claims:       claims.Login{UserID: userID},
						CategoryRule: existingRule,
					}); err != nil {
						log.Error().Err(err).Msg("  updating budget category rule")
						return bizzleException.ErrUnexpected{}
					}
				}
				// go to next rule to shnc
				continue nextRuleToSync
			}
		}

		// if execution reaches here then ruleToSync does not yet exist
		log.Info().Msg(fmt.Sprintf("  Rule %s does not yet exist - Creating it", ruleToSync.Name))

		//create it
		ruleToSync.OwnerID = findOneUserResponse.User.ID
		if _, err := budgetEntryCategoryRuleAdminImp.CreateOne(&budgetEntryCategoryRuleAdmin.CreateOneRequest{
			Claims:       claims.Login{UserID: userID},
			CategoryRule: ruleToSync,
		}); err != nil {
			log.Error().Err(err).Msg("creating budget category rule")
			return bizzleException.ErrUnexpected{}
		}
	}

	log.Info().Msg("SyncBudgetCategoryRulesForUser Run Successfully")
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
