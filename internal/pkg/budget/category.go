package budget

import "strings"

type Category string

const (
	OtherCategory       Category = "Other"
	ElectricityCategory Category = "Electricity"
)

type CategoryIdentifier string

func (c CategoryIdentifier) String() string {
	return string(c)
}

type CategorisationRule struct {
	CategoryIdentifiers []CategoryIdentifier
	Category            Category
	Strict              bool
}

func Categorise(description string) (Category, []CategoryIdentifier, error) {

	// minimise and strip description
	description = strings.ToLower(strings.Trim(description, " "))

	// cannot categorise blank description
	if description == "" {
		return "", nil, ErrCouldNotClassify{Reason: "blank description"}
	}

nextCategorisationRule:
	for _, rule := range CategorisationRules {
		if rule.Strict {
			// all identifiers must be found in description
			for _, id := range rule.CategoryIdentifiers {
				if !strings.Contains(description, id.String()) {
					// if any 1 is not found, go to next rule
					continue nextCategorisationRule
				}
			}
			// if execution reaches here then all category identifiers were found
			return rule.Category, rule.CategoryIdentifiers, nil
		} else {
			// any identifiers can be found in description
			matchedIdentifiers := make([]CategoryIdentifier, 0)
			for _, id := range rule.CategoryIdentifiers {
				if strings.Contains(description, id.String()) {
					// mark that one was found
					matchedIdentifiers = append(matchedIdentifiers, id)
				}
			}
			if len(matchedIdentifiers) > 0 {
				return rule.Category, matchedIdentifiers, nil
			}
		}
	}

	return "", nil, ErrCouldNotClassify{Reason: "not match"}
}

var CategorisationRules = []CategorisationRule{
	{
		CategoryIdentifiers: []CategoryIdentifier{
			"electricity",
			"fee",
		},
		Category: ElectricityCategory,
	},
}
