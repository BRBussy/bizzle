package store

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	CreateMany(CreateManyRequest) (*CreateManyResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
	FindMany(FindManyRequest) (*FindManyResponse, error)
	FindManyComposite(FindManyCompositeRequest) (*FindManyCompositeResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
	DeleteOne(DeleteOneRequest) (*DeleteOneResponse, error)
}

const ServiceProvider = "BudgetEntry-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const FindManyCompositeService = ServiceProvider + ".FindManyComposite"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Entry budgetEntry.Entry
}

type CreateOneResponse struct {
}

type CreateManyRequest struct {
	Entries []budgetEntry.Entry
}

type CreateManyResponse struct {
}

type FindOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

type FindOneResponse struct {
	Entry budgetEntry.Entry
}

type FindManyRequest struct {
	Claims   claims.Claims     `validate:"required"`
	Criteria criteria.Criteria `validate:"required"`
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []budgetEntry.Entry
	Total   int64
}

type FindManyCompositeRequest struct {
	Claims   claims.Claims     `validate:"required"`
	Criteria criteria.Criteria `validate:"required"`
	Query    mongo.Query
}

type FindManyCompositeResponse struct {
	Records []budgetEntry.CompositeEntry
	Total   int64
}

type UpdateOneRequest struct {
	Entry  budgetEntry.Entry
	Claims claims.Claims `validate:"required"`
}

type UpdateOneResponse struct {
}

type DeleteOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

type DeleteOneResponse struct {
}
