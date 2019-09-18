package role

import (
	securityPermission "github.com/BRBussy/bizzle/internal/pkg/security/permission"
	securityRole "github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
)

var initialRoles = []securityRole.Role{
	{
		Name:        "user",
		Permissions: []securityPermission.Permission{},
	},
}

func Setup(
	roleStore roleStore.Store,
) error {

}
