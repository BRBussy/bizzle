package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/user"
)

type Admin interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
	UpdateOne(*UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "User-Store"

const CreateOneService = ServiceProvider + ".CreateOne"
const UpdateOneService = ServiceProvider + ".UpdateOne"

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
