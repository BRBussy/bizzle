package basic

import (
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	tokenGenerator "github.com/BRBussy/bizzle/internal/pkg/security/token/generator"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type authenticator struct {
	userStore        userStore.Store
	tokenGenerator   tokenGenerator.Generator
	requestValidator validationValidator.Validator
}

func New(
	userStore userStore.Store,
	tokenGenerator tokenGenerator.Generator,
	requestValidator validationValidator.Validator,
) bizzleAuthenticator.Authenticator {
	return &authenticator{
		requestValidator: requestValidator,
		userStore:        userStore,
		tokenGenerator:   tokenGenerator,
	}
}

func (a *authenticator) Login(request *bizzleAuthenticator.LoginRequest) (*bizzleAuthenticator.LoginResponse, error) {
	if err := a.requestValidator.ValidateRequest(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and retrieve user by email address
	retrieveResponse, err := a.userStore.FindOne(&userStore.FindOneRequest{
		Identifier: identifier.Email(request.Email),
	})
	if err != nil {
		log.Error().Err(err).Msg("retrieving user for log in")
		return nil, err
	}

	// check password is correct
	if err := bcrypt.CompareHashAndPassword(retrieveResponse.User.Password, []byte(request.Password)); err != nil {
		return nil, err
	}

	// generate login claims
	generateTokenResponse, err := a.tokenGenerator.GenerateToken(
		&tokenGenerator.GenerateTokenRequest{
			Claims: claims.Login{
				UserID:         retrieveResponse.User.ID,
				ExpirationTime: time.Now().Add(time.Hour * 1).UTC().Unix(),
			},
		},
	)
	if err != nil {
		log.Error().Err(err).Msg("generating token")
		return nil, bizzleException.ErrUnexpected{}
	}

	return &bizzleAuthenticator.LoginResponse{
		JWT: generateTokenResponse.Token,
	}, nil
}

func (a *authenticator) AuthenticateService(request *bizzleAuthenticator.AuthenticateServiceRequest) (*bizzleAuthenticator.AuthenticateServiceResponse, error) {
	if err := a.requestValidator.ValidateRequest(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &bizzleAuthenticator.AuthenticateServiceResponse{}, nil
}
