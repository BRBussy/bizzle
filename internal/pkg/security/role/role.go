package role

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/permission"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Role struct {
	ID          identifier.ID           `json:"id" bson:"id"`
	Name        string                  `json:"name" bson:"name"`
	Permissions []permission.Permission `json:"permissions" bson:"permissions"`
}

func (r *Role) AddUniquePermission(p permission.Permission) {
	for pI := range r.Permissions {
		// if role already has permission, return
		if r.Permissions[pI] == p {
			return
		}
	}
	// if execution reaches here then add permission
	r.Permissions = append(r.Permissions, p)
}

func CompareRoles(r1, r2 Role) bool {
	if r1.ID != r2.ID {
		return false
	}
	if r1.Name != r2.Name {
		return false
	}
	if len(r1.Permissions) != len(r2.Permissions) {
		return false
	}
	// for every permission in r1
nextR1Perm:
	for r1PermI := range r1.Permissions {
		// look for it in r2
		for r2PermJ := range r2.Permissions {
			if r1.Permissions[r1PermI] == r2.Permissions[r2PermJ] {
				// if it is found, go to next r1 perm
				continue nextR1Perm
			}
		}
		// if execution reaches here then r1PermI was not found in r2
		return false
	}
	// if execution reaches here every perm in r1 was found in r2
	return true
}
