package role

import (
	securityPermission "github.com/BRBussy/bizzle/internal/pkg/security/permission"
	securityRole "github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin"
	"github.com/rs/zerolog/log"
)

var initialRoles = []securityRole.Role{
	{
		Name:        "user",
		Permissions: []securityPermission.Permission{},
	},
}

func Setup(
	admin roleAdmin.Admin,
) error {
	for i := range initialRoles {
		if _, err := admin.CreateOne(&roleAdmin.CreateOneRequest{Role: initialRoles[i]}); err != nil {
			log.Error().Err(err).Msg("creating role")
			return err
		}
	}
	return nil
}
