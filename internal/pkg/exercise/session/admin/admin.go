package admin

import "github.com/BRBussy/bizzle/internal/pkg/exercise/session"

type Admin interface {
	CreateOne(*CreateOneRequest) (*CreateOneResponse, error)
}

type CreateOneRequest struct {
	Session session.Session `validate:"required"`
}

type CreateOneResponse struct {
	Session session.Session
}