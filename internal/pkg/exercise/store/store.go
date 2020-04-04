package store

import (
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
	FindMany(FindManyRequest) (*FindManyResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "Exercise-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Exercise exercise.Exercise
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Identifier identifier.Identifier
}

type FindOneResponse struct {
	Exercise exercise.Exercise
}

type FindManyRequest struct {
	Criteria criteria.Criteria
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []exercise.Exercise
	Total   int64
}

type UpdateOneRequest struct {
	Exercise exercise.Exercise
}

type UpdateOneResponse struct {
}
