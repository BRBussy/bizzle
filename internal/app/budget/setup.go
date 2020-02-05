package budget

import budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"

var budgetCategoryRulesToMake = []budgetEntryCategoryRule.CategoryRule{
	{
		CategoryIdentifiers: []string{
			"electricity",
			"fee",
		},
		Category: "Electricity",
		Strict:   true,
	},
	{
		CategoryIdentifiers: []string{
			"wesbank",
		},
		Category: "CarRepayment",
	},
	{
		CategoryIdentifiers: []string{
			"vod",
			"prepaid",
		},
		Category: "CellphoneAirtimeData",
		Strict:   true,
	},
	{
		CategoryIdentifiers: []string{
			"telkommobi",
		},
		Category: "Internet",
	},
	{
		CategoryIdentifiers: []string{
			"disc",
			"prem",
			"medical",
		},
		Category: "MedicalAid",
		Strict:   true,
	},
	{
		CategoryIdentifiers: []string{
			"salary",
			"andile",
		},
		Category: "Salary",
		Strict:   true,
	},
}
