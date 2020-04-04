package basic

import (
	roleAdmin "github.com/BRBussy/bizzle/internal/pkg/security/role/admin"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
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
	request.Role.ID = identifier.ID(uuid.NewV4().String())

	if _, err := a.roleStore.CreateOne(roleStore.CreateOneRequest{Role: request.Role}); err != nil {
		log.Error().Err(err).Msg("creating role")
		return nil, err
	}

	return &roleAdmin.CreateOneResponse{Role: request.Role}, nil
}

func (a *admin) UpdateOne(request *roleAdmin.UpdateOneRequest) (*roleAdmin.UpdateOneResponse, error) {
	// try and retrieve the role to be updated
	findOneResponse, err := a.roleStore.FindOne(roleStore.FindOneRequest{Identifier: identifier.ID(request.Role.ID)})
	if err != nil {
		log.Error().Err(err).Msg("finding role to update")
		return nil, err
	}

	// update allowed fields
	findOneResponse.Role.Permissions = request.Role.Permissions

	// update role record
	if _, err := a.roleStore.UpdateOne(roleStore.UpdateOneRequest{Role: findOneResponse.Role}); err != nil {
		log.Error().Err(err).Msg("updating role record")
		return nil, err
	}

	return &roleAdmin.UpdateOneResponse{
		Role: findOneResponse.Role,
	}, nil
}
