package store

import (
	"github.com/BRBussy/bizzle/internal/pkg/user"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Store interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	FindOne(*FindOneRequest) (*FindOneResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "User-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const FindOneService = ServiceProvider + ".FindOne"

type CreateOneRequest struct {
	User user.User
}

type CreateOneResponse struct {
}

type FindOneRequest struct {
	Identifier identifier.Identifier
}

type FindOneResponse struct {
	User user.User
}

type UpdateOneRequest struct {
	User user.User
}

type UpdateOneResponse struct {
}
