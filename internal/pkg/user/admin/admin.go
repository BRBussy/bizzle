package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Admin interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
	RegisterOne(RegisterOneRequest) (*RegisterOneResponse, error)
}

const ServiceProvider = "User-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const UpdateOneService = ServiceProvider + ".UpdateOne"
const RegisterOneService = ServiceProvider + ".RegisterOne"

type CreateOneRequest struct {
	Claims claims.Claims `validate:"required"`
	User   user.User     `validate:"required"`
}

type CreateOneResponse struct {
}

type UpdateOneRequest struct {
	Claims claims.Claims `validate:"required"`
	User   user.User     `validate:"required"`
}

type UpdateOneResponse struct {
}

type RegisterOneRequest struct {
	Claims     claims.Claims         `validate:"required"`
	Identifier identifier.Identifier `validate:"required"`
	Password   string                `validate:"required"`
}

type RegisterOneResponse struct {
}
