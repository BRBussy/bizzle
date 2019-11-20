package admin

import "github.com/BRBussy/bizzle/internal/pkg/exercise/session"

type Admin interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
}

const ServiceProvider = "Session-Admin"

const CreateOneService = ServiceProvider + ".CreateOne"

type CreateOneRequest struct {
	Session session.Session `validate:"required"`
}

type CreateOneResponse struct {
	Session session.Session
}
