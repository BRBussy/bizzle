package user

import (
	goFirebaseAuth "firebase.google.com/go/auth"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/firebase"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
)

var rootUserTemplate = user.User{
	Name:         "root",
	EmailAddress: "root@bizzle.com",
	RoleIDs: []identifier.ID{
		"root",
	},
}

func Setup(
	userAdminImp userAdmin.Admin,
	userStoreImp userStore.Store,
	firebaseImp *firebase.Firebase,
	rootPassword string,
) error {
	// try and retrieve root user's firebase user
	rootUser, err := firebaseImp.GetUserByEmail(rootUserTemplate.EmailAddress)
	if err == nil {
		// user could be retrieved, confirm user's details are correct
		if rootUser.Email != rootUserTemplate.EmailAddress {
			err := bizzleException.ErrUnexpected{Reasons: []string{"root firebase user email incorrect"}}
			log.Error().Err(err)
			return err
		}
		if rootUser.DisplayName != rootUserTemplate.Name {
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
			DisplayName(rootUserTemplate.Name).
			Email(rootUserTemplate.EmailAddress).
			EmailVerified(true).
			Password(rootPassword)
		createdRootUser, err := firebaseImp.CreateUser(firebaseUserToCreateParams)
		if err != nil {
			log.Error().Err(err).Msg("creating root firebase user")
			return bizzleException.ErrUnexpected{}
		}
		rootUser = createdRootUser
	}

	// root user's firebase user retrieved or created
	// try and retrieve root user's bizzle user
	bizzleRootUser, err := userStoreImp.FindOne(userStore.FindOneRequest{
		Identifier: identifier.
			})

	return nil
}
