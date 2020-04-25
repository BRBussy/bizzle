package basic

import (
	"errors"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	userAdmin "github.com/BRBussy/bizzle/internal/pkg/user/admin"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	userValidator "github.com/BRBussy/bizzle/internal/pkg/user/validator"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type admin struct {
	requestValidator validationValidator.Validator
	roleStore        roleStore.Store
	userValidator    userValidator.Validator
	userStore        userStore.Store
}

func New(
	requestValidator validationValidator.Validator,
	userValidator userValidator.Validator,
	userStore userStore.Store,
	roleStore roleStore.Store,
) userAdmin.Admin {
	return &admin{
		requestValidator: requestValidator,
		roleStore:        roleStore,
		userValidator:    userValidator,
		userStore:        userStore,
	}
}

func (a *admin) CreateOne(request userAdmin.CreateOneRequest) (*userAdmin.CreateOneResponse, error) {
	if err := a.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// set user ID
	request.User.ID = identifier.ID(uuid.NewV4().String())

	// validate user for creation
	validateResponse, err := a.userValidator.ValidateForCreate(
		userValidator.ValidateForCreateRequest{
			User: request.User,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("error validating user for create")
		return nil, bizzleException.ErrUnexpected{}
	}
	if len(validateResponse.ReasonsInvalid) > 0 {
		return nil, userAdmin.ErrUserInvalid{ReasonsInvalid: validateResponse.ReasonsInvalid}
	}

	// create user database entity
	if _, err := a.userStore.CreateOne(
		userStore.CreateOneRequest{
			User: request.User,
		},
	); err != nil {
		log.Error().Err(err).Msg("error creating user")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &userAdmin.CreateOneResponse{}, nil
}

func (a *admin) UpdateOne(userAdmin.UpdateOneRequest) (*userAdmin.UpdateOneResponse, error) {
	return nil, errors.New("implement update one")
}

func (a *admin) RegisterOne(request userAdmin.RegisterOneRequest) (*userAdmin.RegisterOneResponse, error) {
	if err := a.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// retrieve user being registered
	retrieveUserResponse, err := a.userStore.FindOne(
		userStore.FindOneRequest{
			Identifier: request.UserIdentifier,
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("error finding user to register")
		return nil, bizzleException.ErrUnexpected{}
	}

	// hash user password
	pwdHash, err := bcrypt.GenerateFromPassword(
		[]byte(request.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		log.Error().Err(err).Msg("error hashing user's password")
		return nil, bizzleException.ErrUnexpected{}
	}

	// set user password
	retrieveUserResponse.User.Password = pwdHash

	// set user to registered
	retrieveUserResponse.User.Registered = true

	// update user
	if _, err := a.userStore.UpdateOne(
		userStore.UpdateOneRequest{
			User: retrieveUserResponse.User,
		},
	); err != nil {
		log.Error().Err(err).Msg("error updating user")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &userAdmin.RegisterOneResponse{}, nil
}
