package basic

import (
	sessionAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/session/admin"
	sessionStore "github.com/BRBussy/bizzle/internal/pkg/exercise/session/store"
	validateValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type admin struct {
	sessionStore     sessionStore.Store
	requestValidator validateValidator.Validator
}

func New(
	sessionStore sessionStore.Store,
	requestValidator validateValidator.Validator,
) sessionAdmin.Admin {
	return &admin{
		requestValidator: requestValidator,
		sessionStore:     sessionStore,
	}
}

func (a admin) CreateOne(request *sessionAdmin.CreateOneRequest) (*sessionAdmin.CreateOneResponse, error) {
	if err := a.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &sessionAdmin.CreateOneResponse{}, nil
}
