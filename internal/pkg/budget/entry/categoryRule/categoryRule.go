package categoryRule

import (
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type CategoryRule struct {
	ID                  identifier.ID
	OwnerID             identifier.ID
	CategoryIdentifiers []string
	Category            string
	Strict              bool
}
