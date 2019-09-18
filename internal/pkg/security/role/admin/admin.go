package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
)

type Admin interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
}

const ServiceProvider = "Role-Admin"

const CreateOneService = ServiceProvider + ".CreateOne"

type CreateOneRequest struct {
	Role role.Role
}

type CreateOneResponse struct {
	Role role.Role
}
