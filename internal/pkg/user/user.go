package user

import "github.com/BRBussy/bizzle/pkg/search/identifier"

type User struct {
	ID          string          `json:"id" bson:"id"`
	Name        string          `json:"name" bson:"name"`
	Email       string          `json:"email" bson:"email"`
	RoleIDs     []identifier.ID `json:"roleIDs" bson:"roleIDs"`
	FirebaseUID string          `json:"firebaseUID" bson:"firebaseUID"`
}
