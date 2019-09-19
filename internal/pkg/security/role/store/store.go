package store

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	FindOne(*FindOneRequest) (*FindOneResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "Role-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Role role.Role
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Identifier identifier.Identifier
}

type FindOneResponse struct {
	Role role.Role
}

type UpdateOneRequest struct {
	Role role.Role
}

type UpdateOneResponse struct {
}
