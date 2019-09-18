package basic

import (
	roleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/rs/zerolog/log"
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

	if _, err := a.roleStore.CreateOne(&roleStore.CreateOneRequest{Role: request.Role}); err != nil {
		log.Error().Err(err).Msg("creating role")
		return nil, err
	}

	return &roleAdmin.CreateOneResponse{Role: request.Role}, nil
}
