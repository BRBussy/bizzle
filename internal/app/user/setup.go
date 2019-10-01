package user

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	stringCriterion "github.com/BRBussy/bizzle/pkg/search/criterion/string"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

var rootUserToCreate = user.User{
	Name:  "root",
	Email: "root@bizzle.com",
	RoleIDs: []identifier.ID{
		"root",
	},
}

func Setup(
	userAdminImp userAdmin.Admin,
	userStoreImp userStore.Store,
	roleStoreImp roleStore.Store,
	rootPassword string,
) error {
	// if root user to be created has any roles, find the role ids now
	if len(rootUserToCreate.RoleIDs) > 0 {
		roleFindCriteria := make([]criterion.Criterion, 0)
		for i := range rootUserToCreate.RoleIDs {
			roleFindCriteria = append(
				roleFindCriteria,
				stringCriterion.Exact{
					Field:  "name",
					String: rootUserToCreate.RoleIDs[i].String(),
				},
			)
		}
		roleFindResponse, err := roleStoreImp.FindMany(&roleStore.FindManyRequest{
			Criteria: roleFindCriteria,
		})
		if err != nil {
			log.Error().Err(err).Msg("finding root user roles")
			return bizzleException.ErrUnexpected{}
		}
		rootUserToCreate.RoleIDs = make([]identifier.ID, 0)
		for i := range roleFindResponse.Records {
			rootUserToCreate.RoleIDs = append(rootUserToCreate.RoleIDs, roleFindResponse.Records[i].ID)
		}
	}

	// root user's firebase user retrieved or created
	// try and retrieve root user's bizzle user
	findRootUserResponse, err := userStoreImp.FindOne(&userStore.FindOneRequest{
		Identifier: identifier.Email(rootUserToCreate.Email),
	})
	if err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			// root user in bizzle not found, create it
			createResponse, err := userAdminImp.CreateOne(&userAdmin.CreateOneRequest{User: rootUserToCreate})
			if err != nil {
				log.Error().Err(err).Msg("creating bizzle root user")
				return bizzleException.ErrUnexpected{}
			}
			findRootUserResponse = &userStore.FindOneResponse{User: createResponse.User}
		default:
			log.Error().Err(err).Msg("retrieving root user")
			return bizzleException.ErrUnexpected{}
		}
	}

	// root user has been found or created, populate id before comparison
	rootUserToCreate.ID = findRootUserResponse.User.ID

	// check if found user is different from that which should be created
	if !user.CompareUsers(rootUserToCreate, findRootUserResponse.User) {
		return bizzleException.ErrUnexpected{Reasons: []string{"root user not in sync"}}
	}

	return nil
}
