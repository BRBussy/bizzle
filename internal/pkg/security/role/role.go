package role

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/permission"
)

type Role struct {
	ID          string                  `json:"id" bson:"id"`
	Name        string                  `json:"name" bson:"name"`
	Permissions []permission.Permission `json:"permissions" bson:"permissions"`
}
