package store

import (
	budgetIgnoredIgnored "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored"
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
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
	DeleteOne(DeleteOneRequest) (*DeleteOneResponse, error)
}

const ServiceProvider = "BudgetIgnoredIgnored-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Ignored budgetIgnoredIgnored.Ignored
}

type CreateOneResponse struct {
}

type CreateManyRequest struct {
	Entries []budgetIgnoredIgnored.Ignored
}

type CreateManyResponse struct {
}

type FindOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

type FindOneResponse struct {
	Ignored budgetIgnoredIgnored.Ignored
}

type FindManyRequest struct {
	Claims   claims.Claims     `validate:"required"`
	Criteria criteria.Criteria `validate:"required"`
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []budgetIgnoredIgnored.Ignored
	Total   int64
}

type UpdateOneRequest struct {
	Ignored budgetIgnoredIgnored.Ignored
	Claims  claims.Claims `validate:"required"`
}

type UpdateOneResponse struct {
}

type DeleteOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

type DeleteOneResponse struct {
}
