package entry

import (
	"github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"time"
)

type Entry struct {
	ID             identifier.ID `json:"id" bson:"id"`
	OwnerID        identifier.ID `json:"ownerID" bson:"ownerID"`
	Date           time.Time     `json:"date" bson:"date"`
	Description    string        `json:"description" bson:"description"`
	Amount         float64       `json:"amount" bson:"amount"`
	CategoryRuleID identifier.ID `json:"categoryRuleID" bson:"categoryRuleID"`
}

type CompositeEntry struct {
	Entry        `bson:"inline"`
	CategoryRule categoryRule.CategoryRule `json:"categoryRule"`
}
