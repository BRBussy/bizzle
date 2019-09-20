package user

import (
	goFirebaseAuth "firebase.google.com/go/auth"
	"fmt"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/firebase"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
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
	firebaseImp *firebase.Firebase,
	rootPassword string,
) error {
	// try and retrieve root user's firebase user
	rootUser, err := firebaseImp.GetUserByEmail(rootUserToCreate.Email)
	if err == nil {
		// user could be retrieved, confirm user's details are correct
		if rootUser.Email != rootUserToCreate.Email {
			err := bizzleException.ErrUnexpected{Reasons: []string{"root firebase user email incorrect"}}
			log.Error().Err(err)
			return err
		}
		if rootUser.DisplayName != rootUserToCreate.Name {
			err := bizzleException.ErrUnexpected{Reasons: []string{"root firebase user name incorrect"}}
			log.Error().Err(err)
			return err
		}
		if !rootUser.EmailVerified {
			err := bizzleException.ErrUnexpected{Reasons: []string{"root firebase user not email verified"}}
			log.Error().Err(err)
			return err
		}
	} else {
		// user could not be retrieved, try and create user
		firebaseUserToCreateParams := (&goFirebaseAuth.UserToCreate{}).
			DisplayName(rootUserToCreate.Name).
			Email(rootUserToCreate.Email).
			EmailVerified(true).
			Password(rootPassword)
		createdRootUser, err := firebaseImp.CreateUser(firebaseUserToCreateParams)
		if err != nil {
			log.Error().Err(err).Msg("creating root firebase user")
			return bizzleException.ErrUnexpected{}
		}
		rootUser = createdRootUser
	}
	// update firebase ID on root user to create
	rootUserToCreate.FirebaseUID = rootUser.UID

	// root user's firebase user retrieved or created
	// try and retrieve root user's bizzle user
	findRootUserResponse, err := userStoreImp.FindOne(&userStore.FindOneRequest{
		Identifier: identifier.Email(rootUserToCreate.Email),
	})
	if err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			// root user in bizzle not found, create it
			// first retrieve the roles for root user
			roleFindResponse, err := roleStoreImp.FindMany(&roleStore.FindManyRequest{})
			if err != nil {
				log.Error().Err(err).Msg("finding root user roles")
				return bizzleException.ErrUnexpected{}
			}
			fmt.Println("find all roles!", roleFindResponse)

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

	// root user has been found or created
	fmt.Println(findRootUserResponse)
	return nil
}
