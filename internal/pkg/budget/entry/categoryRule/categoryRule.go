package categoryRule

import (
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type CategoryRule struct {
	ID                   identifier.ID `json:"id" bson:"id"`
	OwnerID              identifier.ID `validate:"required" json:"ownerID" bson:"ownerID"`
	CategoryIdentifiers  []string      `validate:"required,gt=0" json:"categoryIdentifiers" bson:"categoryIdentifiers"`
	Name                 string        `validate:"required" json:"name" bson:"name"`
	Strict               bool          `json:"strict" bson:"strict"`
	ExpectedAmount       float32       `json:"expectedAmount" bson:"expectedAmount"`
	ExpectedAmountPeriod int           `json:"expectedAmountPeriod" bson:"expectedAmountPeriod"`
}

func CompareCategoryRules(c1, c2 CategoryRule) bool {
	if c1.ID != c2.ID {
		return false
	}
	if c1.OwnerID != c2.OwnerID {
		return false
	}
	if c1.Name != c2.Name {
		return false
	}
	if c1.Strict != c2.Strict {
		return false
	}
	if c1.ExpectedAmount != c2.ExpectedAmount {
		return false
	}
	if c1.ExpectedAmountPeriod != c2.ExpectedAmountPeriod {
		return false
	}
	if len(c1.CategoryIdentifiers) != len(c2.CategoryIdentifiers) {
		return false
	}
	// for every category in c1
nextC1Cat:
	for c1CatI := range c1.CategoryIdentifiers {
		// look for it in c2
		for c2CatJ := range c2.CategoryIdentifiers {
			if c1.CategoryIdentifiers[c1CatI] == c2.CategoryIdentifiers[c2CatJ] {
				// if it is found, go to next c1 cat
				continue nextC1Cat
			}
		}
		// if execution reaches here then c1CatI was not found in c2
		return false
	}
	// if execution reaches here then every cat in c1 was found in c2
	return true
}
