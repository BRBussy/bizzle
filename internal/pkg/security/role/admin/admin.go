package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
)

type Admin interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	UpdateOne(UpdateOneRequest) (*UpdateOneResponse, error)
}

const ServiceProvider = "Role-Admin"

const CreateOneService = ServiceProvider + ".CreateOne"
const UpdateOneService = ServiceProvider + ".UpdateOne"

type CreateOneRequest struct {
	Role role.Role
}

type CreateOneResponse struct {
	Role role.Role
}

type UpdateOneRequest struct {
	Role role.Role
}

type UpdateOneResponse struct {
	Role role.Role
}
