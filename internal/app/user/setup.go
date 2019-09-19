package user

import (
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
)

func Setup(
	userAdmin userAdmin.Admin,
	userStore userStore.Store,
) error {
	return nil
}
