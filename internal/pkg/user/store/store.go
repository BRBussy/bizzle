package store

import (
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
	FindMany(FindManyRequest) (*FindManyResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "User-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	User user.User
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

type FindOneResponse struct {
	User user.User
}

type FindManyRequest struct {
	Claims   claims.Claims     `validate:"required"`
	Criteria criteria.Criteria `validate:"required"`
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []user.User
	Total   int64
}

type UpdateOneRequest struct {
	Claims claims.Claims `validate:"required"`
	User   user.User
}

type UpdateOneResponse struct {
}
