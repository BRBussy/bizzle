package entry

import (
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Entry struct {
	ID          identifier.ID `json:"id" bson:"id"`
	OwnerID     identifier.ID `json:"ownerID" bson:"ownerID"`
	Date        string        `json:"date" bson:"date"`
	Description string        `json:"description" bson:"description"`
	Amount      float64       `json:"amount" bson:"amount"`
	Category    Category      `json:"category" bson:"category"`
}
