package user

import (
	"github.com/BRBussy/bizzle/internal/pkg/firebase"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
)

func Setup(
	userAdmin userAdmin.Admin,
	userStore userStore.Store,
	firebase *firebase.Firebase,
) error {

	return nil
}
