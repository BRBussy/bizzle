package store

import "github.com/BRBussy/bizzle/internal/pkg/security/role"

type Store interface {
	Create(*CreateRequest) (*CreateResponse, error)
}

const ServiceProvider = "Role-Store"

const CreateService = ServiceProvider + ".Create"

type CreateRequest struct {
	Role role.Role
}

type CreateResponse struct {
	Role role.Role
}
