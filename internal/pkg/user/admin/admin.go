package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
)

type Admin interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	RegisterOne(RegisterOneRequest) (*RegisterOneResponse, error)
	ChangePassword(ChangePasswordRequest) (*ChangePasswordResponse, error)
}

const ServiceProvider = "User-Admin"

const CreateOneService = ServiceProvider + ".CreateOne"
const RegisterOneService = ServiceProvider + ".RegisterOne"
const ChangePasswordService = ServiceProvider + ".ChangePassword"

type CreateOneRequest struct {
	Claims claims.Claims `validate:"required"`
	User   user.User     `validate:"required"`
}

type CreateOneResponse struct {
}

type RegisterOneRequest struct {
	Claims         claims.Claims         `validate:"required"`
	UserIdentifier identifier.Identifier `validate:"required"`
	Password       string                `validate:"required"`
}

type RegisterOneResponse struct {
}

type ChangePasswordRequest struct {
	Claims         claims.Claims         `validate:"required"`
	UserIdentifier identifier.Identifier `validate:"required"`
	Password       string                `validate:"required"`
}

type ChangePasswordResponse struct {
}
