package entry

import (
	"time"

	"github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

// Entry is a budget entry
type Entry struct {
	ID             identifier.ID `json:"id" bson:"id"`
	OwnerID        identifier.ID `validate:"required" json:"ownerID" bson:"ownerID"`
	Date           time.Time     `validate:"required,date" json:"date" bson:"date"`
	Description    string        `validate:"required" json:"description" bson:"description"`
	Amount         float64       `validate:"required" json:"amount" bson:"amount"`
	CategoryRuleID identifier.ID `json:"categoryRuleID" bson:"categoryRuleID"`
}

// CompositeEntry is used to populate budget entry with corresponding category rule
type CompositeEntry struct {
	Entry        `bson:"inline"`
	CategoryRule categoryRule.CategoryRule `json:"categoryRule"`
}
