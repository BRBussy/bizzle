package store

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	Create(*CreateRequest) (*CreateResponse, error)
	FindOne(*FindOneRequest) (*FindOneResponse, error)
}

const ServiceProvider = "Role-Store"

const CreateService = ServiceProvider + ".Create"
const FindOneService = ServiceProvider + ".FindOne"

type CreateRequest struct {
	Role role.Role
}

type CreateResponse struct {
	Role role.Role
}

type FindOneRequest struct {
	Identifier identifier.Identifier
}

type FindOneResponse struct {
	Role role.Role
}
