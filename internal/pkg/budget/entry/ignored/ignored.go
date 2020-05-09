package ignored

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type Ignored struct {
	ID          identifier.ID `json:"id" bson:"id"`
	OwnerID     identifier.ID `json:"ownerID" bson:"ownerID"`
	Description string        `validate:"required" json:"description" bson:"description"`
}
