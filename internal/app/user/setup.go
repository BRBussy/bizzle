package user

import (
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
		log.Error().Err(err).Msg("getting root user by email")
		return bizzleException.ErrUnexpected{}
	}

	return nil
}
