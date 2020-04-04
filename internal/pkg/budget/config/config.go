package config

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type Config struct {
	ID                  identifier.ID `json:"id" bson:"id"`
	OwnerID             identifier.ID `validate:"required" json:"ownerID" bson:"ownerID"`
	OtherCategoryRuleID identifier.ID `json:"otherCategoryRuleID" bson:"otherCategoryRuleID"`
}
