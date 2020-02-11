package store

import (
	budgetEntryCategoryRule "github.com/BRBussy/bizzle/internal/pkg/budget/entry/categoryRule"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

// Store is used to perform crud operations on budget category rules
type Store interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	CreateMany(*CreateManyRequest) (*CreateManyResponse, error)
	FindOne(*FindOneRequest) (*FindOneResponse, error)
	FindMany(*FindManyRequest) (*FindManyResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "BudgetEntryCategoryRule-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	CategoryRule budgetEntryCategoryRule.CategoryRule
}

type CreateOneResponse struct {
}

type CreateManyRequest struct {
	Entries []budgetEntryCategoryRule.CategoryRule
}

type CreateManyResponse struct {
}

type FindOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

type FindOneResponse struct {
	CategoryRule budgetEntryCategoryRule.CategoryRule
}

type FindManyRequest struct {
	Criteria criteria.Criteria `validate:"required"`
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []budgetEntryCategoryRule.CategoryRule
	Total   int64
}

type UpdateOneRequest struct {
	CategoryRule budgetEntryCategoryRule.CategoryRule
}

type UpdateOneResponse struct {
}
