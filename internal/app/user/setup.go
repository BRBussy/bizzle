package user

import (
	goFirebaseAuth "firebase.google.com/go/auth"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/firebase"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/rs/zerolog/log"
)

func Setup(
	userAdmin userAdmin.Admin,
	userStore userStore.Store,
	firebase *firebase.Firebase,
) error {
	// try and retrieve root firebase user
	rootUser, err := firebase.GetUserByEmail("root@bizzle.com")
	if err != nil {
		// user could not be retrieved, try and create user
		firebaseUserToCreateParams := (&goFirebaseAuth.UserToCreate{}).
			DisplayName("root").
			Email("root@bizzle.com").
			EmailVerified(true).
			Password("123456")
		createdRootUser, err := firebase.CreateUser(firebaseUserToCreateParams)
		if err != nil {
			log.Error().Err(err).Msg("creating root firebase user")
			return bizzleException.ErrUnexpected{}
		}
		rootUser = createdRootUser
	}

	return nil
}
