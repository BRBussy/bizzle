package store

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/BRBussy/bizzle/pkg/search/query"
)

type Store interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	FindOne(*FindOneRequest) (*FindOneResponse, error)
	FindMany(*FindManyRequest) (*FindManyResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "Role-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
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

type FindManyRequest struct {
	Criteria criteria.Criteria
	Query    query.Query
}

type FindManyResponse struct {
	Records []role.Role
	Total   int
}

type UpdateOneRequest struct {
	Role role.Role
}

type UpdateOneResponse struct {
}
