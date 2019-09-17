package role

import "github.com/BRBussy/bizzle/internal/pkg/security/permission/api"

type Role struct {
	ID             string           `json:"id" bson:"id"`
	Name           string           `json:"name" bson:"name"`
	APIPermissions []api.Permission `json:"apiPermissions"`
}
