package entry

import (
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"time"
)

type Entry struct {
	ID          identifier.ID `json:"id" bson:"id"`
	OwnerID     identifier.ID `json:"ownerID" bson:"ownerID"`
	Date        time.Time     `json:"date" bson:"date"`
	Description string        `json:"description" bson:"description"`
	Amount      float64       `json:"amount" bson:"amount"`
	CategoryID  identifier.ID `json:"category" bson:"category"`
}
