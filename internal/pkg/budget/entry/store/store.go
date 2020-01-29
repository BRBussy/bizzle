package store

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	CreateMany(*CreateManyRequest) (*CreateManyResponse, error)
	FindOne(*FindOneRequest) (*FindOneResponse, error)
	FindMany(*FindManyRequest) (*FindManyResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "BudgetEntry-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
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
	Identifier identifier.Identifier
}

type FindOneResponse struct {
	Entry budgetEntry.Entry
}

type FindManyRequest struct {
	Criteria criteria.Criteria
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []budgetEntry.Entry
	Total   int64
}

type UpdateOneRequest struct {
	Entry budgetEntry.Entry
}

type UpdateOneResponse struct {
}
