package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/user"
)

type Admin interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
	RegisterOne(*RegisterOneRequest) (*RegisterOneResponse, error)
}

const ServiceProvider = "User-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const UpdateOneService = ServiceProvider + ".UpdateOne"
const RegisterOneService = ServiceProvider + ".RegisterOne"

type CreateOneRequest struct {
	User user.User
}

type CreateOneResponse struct {
	User user.User
}

type UpdateOneRequest struct {
	User user.User
}

type UpdateOneResponse struct {
}

type RegisterOneRequest struct {
	User     user.User
	Password string
}

type RegisterOneResponse struct {
}
