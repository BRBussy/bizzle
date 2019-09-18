package role

import (
	securityPermission "github.com/BRBussy/bizzle/internal/pkg/security/permission"
	securityRole "github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/rs/zerolog/log"
)

var initialRoles = []securityRole.Role{
	{
		Name:        "user",
		Permissions: []securityPermission.Permission{},
	},
}

func Setup(
	store roleStore.Store,
) error {
	for i := range initialRoles {
		if _, err := store.Create(&roleStore.CreateRequest{Role: initialRoles[i]}); err != nil {
			log.Error().Err(err).Msg("creating role")
			return err
		}
	}
	return nil
}
