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

// SyncBudgetCategoryRulesForUser is used to sync budget category rules for a user
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
		Claims:   claims.Login{UserID: userID},
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
			"pre-paid electricity",
			"pre paid electricity",
			"electricity",
		},
		Strict: false,
	},
	{
		Name: "Car Repayment",
		CategoryIdentifiers: []string{
			"wesbank",
		},
	},
	{
		Name: "Cellphone Airtime/Data",
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
		Name: "Medical Aid",
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
	{
		Name: "Petrol",
		CategoryIdentifiers: []string{
			"shell",
			"engen",
			"cotswold energy centre",
			"total",
			"petroport n3",
			"oostermoed garage",
			"bp hillcrest",
		},
		Strict: false,
	},
	{
		Name: "Groceries",
		CategoryIdentifiers: []string{
			"checkers",
			"spar",
			"woolworths",
			"pnp",
			"montagu",
			"superspar",
			"impala fruit & flower",
			"saffron spice box",
			"252 pine avn",
			"the flower & nut marke",
		},
		Strict: false,
	},
	{
		Name: "Eating Out",
		CategoryIdentifiers: []string{
			"the living room",
			"turn n tender",
			"makaranga",
			"lexis eatery",
			"savior cafe",
			"little india",
			"marulas coffee shop",
			"factory on grant",
			"fireroom bedfordview",
			"stones bedfordview",
			"wimpy",
			"fournos",
			"anat",
			"free food diner",
			"food lovers market",
			"simply asia",
			"kauai",
			"pauls homemade",
			"dosa hut",
			"planet fit cafe",
			"tikka n kebab",
			"spur",
			"mugg and bean",
			"uber eats",
			"fego",
			"we make coffee",
		},
		Strict: false,
	},
	{
		Name: "Medicine",
		CategoryIdentifiers: []string{
			"dischem",
			"linksfield pharmac",
			"ferngate pharmacy",
		},
		Strict: false,
	},
	{
		Name: "Rent",
		CategoryIdentifiers: []string{
			"rentjune",
			"richardrent",
		},
		Strict: false,
	},
	{
		Name: "Rehab",
		CategoryIdentifiers: []string{
			"relapse prevention",
		},
		Strict: true,
	},
	{
		Name: "Overdraft",
		CategoryIdentifiers: []string{
			"overdraft service fee",
			"overdraft interest",
		},
		Strict: false,
	},
	{
		Name: "Bank Account Fee",
		CategoryIdentifiers: []string{
			"fixed monthly fee",
		},
		Strict: true,
	},
	{
		Name: "Cash Withdrawal",
		CategoryIdentifiers: []string{
			"cash withdrawal",
			"cash withd",
		},
		Strict: false,
	},
	{
		Name: "Muay Thai",
		CategoryIdentifiers: []string{
			"primal gym",
		},
		Strict: true,
	},
	{
		Name: "Hiking",
		CategoryIdentifiers: []string{
			"jhb hiking club",
		},
		Strict: true,
	},
}
