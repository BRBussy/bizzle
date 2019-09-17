package store

import "github.com/BRBussy/bizzle/internal/pkg/security/role"

type Store interface {
	Create(*CreateRequest) (*CreateResponse, error)
}

const ServiceProvider = "Role-Store"

type CreateRequest struct {
	Role role.Role
}

type CreateResponse struct {
	Role role.Role
}
