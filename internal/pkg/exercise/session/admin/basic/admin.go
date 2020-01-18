package basic

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	sessionAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/session/admin"
	sessionStore "github.com/BRBussy/bizzle/internal/pkg/exercise/session/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validateValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	sessionStore     sessionStore.Store
	requestValidator validateValidator.Validator
}

func New(
	requestValidator validateValidator.Validator,
	sessionStore sessionStore.Store,
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

	request.Session.ID = identifier.ID(uuid.NewV4().String())

	if _, err := a.sessionStore.CreateOne(&sessionStore.CreateOneRequest{Session: request.Session}); err != nil {
		log.Error().Err(err).Msg("creating session")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &sessionAdmin.CreateOneResponse{
		Session: request.Session,
	}, nil
}
