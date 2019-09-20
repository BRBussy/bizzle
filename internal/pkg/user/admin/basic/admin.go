package basic

import (
	"errors"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/firebase"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	userValidator "github.com/BRBussy/bizzle/internal/pkg/user/validator"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	roleStore     roleStore.Store
	firebase      *firebase.Firebase
	userValidator userValidator.Validator
	userStore     userStore.Store
}

func New(
	userValidator userValidator.Validator,
	userStore userStore.Store,
	roleStore roleStore.Store,
	firebase *firebase.Firebase,
) userAdmin.Admin {
	return &admin{
		roleStore:     roleStore,
		firebase:      firebase,
		userValidator: userValidator,
		userStore:     userStore,
	}
}

func (a *admin) CreateOne(request *userAdmin.CreateOneRequest) (*userAdmin.CreateOneResponse, error) {
	validateResponse, err := a.userValidator.ValidateForCreate(&userValidator.ValidateForCreateRequest{User: request.User})
	if err != nil {
		log.Error().Err(err).Msg("validating user for create")
		return nil, bizzleException.ErrUnexpected{}
	}
	if len(validateResponse.ReasonsInvalid) > 0 {
		return nil, userAdmin.ErrUserInvalid{ReasonsInvalid: validateResponse.ReasonsInvalid}
	}

	request.User.ID = identifier.ID(uuid.NewV4().String())

	if _, err := a.userStore.CreateOne(&userStore.CreateOneRequest{User: request.User}); err != nil {
		log.Error().Err(err).Msg("creating user")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &userAdmin.CreateOneResponse{User: request.User}, nil
}

func (a *admin) UpdateOne(*userAdmin.UpdateOneRequest) (*userAdmin.UpdateOneResponse, error) {
	return nil, errors.New("implement update one")
}
