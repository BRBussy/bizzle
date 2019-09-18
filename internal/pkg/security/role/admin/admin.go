package admin

import (
	"github.com/BRBussy/bizzle/internal/pkg/security/role"
)

type Admin interface {
	Create(*CreateRequest) (*CreateResponse, error)
}

const ServiceProvider = "Role-Admin"

const CreateService = ServiceProvider + ".Create"

type CreateRequest struct {
	Role role.Role
}

type CreateResponse struct {
	Role role.Role
}
