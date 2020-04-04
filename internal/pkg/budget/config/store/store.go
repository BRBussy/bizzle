package store

import (
	budgetConfig "github.com/BRBussy/bizzle/internal/pkg/budget/config"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/pkg/search/criteria"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
	FindMany(FindManyRequest) (*FindManyResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
	DeleteOne(DeleteOneRequest) (*DeleteOneResponse, error)
}

const ServiceProvider = "BudgetEntry-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"
const FindManyService = ServiceProvider + ".FindMany"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Config budgetConfig.Config
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
}

type FindOneResponse struct {
	Config budgetConfig.Config
}

type FindManyRequest struct {
	Claims   claims.Claims     `validate:"required"`
	Criteria criteria.Criteria `validate:"required"`
	Query    mongo.Query
}

type FindManyResponse struct {
	Records []budgetConfig.Config
	Total   int64
}

type UpdateOneRequest struct {
	Config budgetConfig.Config
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
