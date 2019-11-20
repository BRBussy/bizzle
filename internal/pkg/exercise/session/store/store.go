package store

import (
	"github.com/BRBussy/bizzle/internal/pkg/exercise/session"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	FindOne(*FindOneRequest) (*FindOneResponse, error)
	FindMany(*FindManyRequest) (*FindManyResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "Session-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Session session.Session
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Identifier identifier.Identifier
}

type FindOneResponse struct {
	Session session.Session
}

type FindManyRequest struct {
	Criteria criteria.Criteria
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []session.Session
	Total   int64
}

type UpdateOneRequest struct {
	Session session.Session
}

type UpdateOneResponse struct {
}
