package role

import (
	"fmt"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	securityPermission "github.com/BRBussy/bizzle/internal/pkg/security/permission"
	securityRole "github.com/BRBussy/bizzle/internal/pkg/security/role"
	roleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
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
	store roleStore.Store,
) error {
	// for every initial role to create
	for i := range initialRoles {
		// try and retrieve the role
		fineOneResponse, err := store.FindOne(&roleStore.FindOneRequest{Identifier: identifier.Name(initialRoles[i].Name)})
		if err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				// role was not found, create it and move on to next role
				if _, err := admin.CreateOne(&roleAdmin.CreateOneRequest{Role: initialRoles[i]}); err != nil {
					log.Error().Err(err).Msg("creating role")
					return err
				}
				continue

			default:
				// there was some error retrieving the role
				log.Error().Err(err).Msg("finding role")
				return err
			}
		}
		fmt.Println("role found!", fineOneResponse)
	}
	return nil
}
