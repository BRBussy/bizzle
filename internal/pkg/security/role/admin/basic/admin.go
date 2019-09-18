package basic

import (
	roleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	roleStore roleStore.Store
}

func New(
	roleStore roleStore.Store,
) roleAdmin.Admin {
	return &admin{
		roleStore: roleStore,
	}
}

func (a *admin) CreateOne(request *roleAdmin.CreateOneRequest) (*roleAdmin.CreateOneResponse, error) {
	request.Role.ID = uuid.NewV4().String()

	createOneResponse, err := a.roleStore.Create(roleStore.CreateRequest{Role: request.Role})
}
