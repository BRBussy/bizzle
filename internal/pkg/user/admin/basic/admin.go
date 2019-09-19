package basic

import (
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
)

type admin struct {
	roleStore roleStore.Store
}

func New(
	roleStore roleStore.Store,
) userAdmin.Admin {
	return &admin{
		roleStore: roleStore,
	}
}

func (a *admin) CreateOne(*userAdmin.CreateOneRequest) (*userAdmin.CreateOneResponse, error) {
	panic("implement me")
}

func (a *admin) UpdateOne(*userAdmin.UpdateOneRequest) (*userAdmin.UpdateOneResponse, error) {
	panic("implement me")
}
