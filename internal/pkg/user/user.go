package user

import (
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type User struct {
	ID          identifier.ID   `json:"id" bson:"id"`
	Name        string          `validate:"required" json:"name" bson:"name"`
	Email       string          `validate:"required,email" json:"email" bson:"email"`
	RoleIDs     []identifier.ID `validate:"required" json:"roleIDs" bson:"roleIDs"`
	FirebaseUID string          `validate:"required" json:"firebaseUID" bson:"firebaseUID"`
}

func CompareUsers(a, b User) bool {
	if a.ID != b.ID {
		return false
	}
	if a.Name != b.Name {
		return false
	}
	if a.Email != b.Email {
		return false
	}
	if len(a.RoleIDs) != len(b.RoleIDs) {
		return false
	}
	// for every role in a
nextRoleIDa:
	for roleAIdx := range a.RoleIDs {
		// look for the role in b
		for roleBIdx := range b.RoleIDs {
			if b.RoleIDs[roleBIdx] == a.RoleIDs[roleAIdx] {
				// a found in b, consider next a
				continue nextRoleIDa
			}
		}
		// if execution reaches here roleA was not found in b
		return false
	}
	if a.FirebaseUID != b.FirebaseUID {
		return false
	}
	return true
}
