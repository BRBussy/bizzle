package basic

import (
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/security/claims"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	tokenGenerator "github.com/BRBussy/bizzle/internal/pkg/security/token/generator"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/criterion"
	"github.com/BRBussy/bizzle/pkg/search/criterion/text"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type authenticator struct {
	userStore        userStore.Store
	roleStore        roleStore.Store
	tokenGenerator   tokenGenerator.Generator
	requestValidator validationValidator.Validator
	database         *mongo.Database
}

func New(
	userStore userStore.Store,
	roleStore roleStore.Store,
	tokenGenerator tokenGenerator.Generator,
	requestValidator validationValidator.Validator,
	database *mongo.Database,
) bizzleAuthenticator.Authenticator {
	return &authenticator{
		roleStore:        roleStore,
		requestValidator: requestValidator,
		userStore:        userStore,
		tokenGenerator:   tokenGenerator,
		database:         database,
	}
}

func (a *authenticator) Login(request bizzleAuthenticator.LoginRequest) (*bizzleAuthenticator.LoginResponse, error) {
	if err := a.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var userLoggingIn user.User
	err := a.database.Collection("user").FindOne(&userLoggingIn, request.Email)
	if err != nil {
		log.Error().Err(err).Msg("retrieving user for log in")
		return nil, err
	}

	// check password is correct
	if err := bcrypt.CompareHashAndPassword(userLoggingIn.Password, []byte(request.Password)); err != nil {
		log.Error().Err(err).Msg("invalid password login")
		return nil, err
	}

	// generate login claims
	generateTokenResponse, err := a.tokenGenerator.GenerateToken(
		&tokenGenerator.GenerateTokenRequest{
			Claims: claims.Login{
				UserID:         userLoggingIn.ID,
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

func (a *authenticator) AuthenticateService(request bizzleAuthenticator.AuthenticateServiceRequest) (*bizzleAuthenticator.AuthenticateServiceResponse, error) {
	if err := a.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	switch typedClaims := request.Claims.(type) {
	case claims.Login:
		// try and retrieve user that owns claims
		findOneUserResponse, err := a.userStore.FindOne(
			userStore.FindOneRequest{
				Claims:     request.Claims,
				Identifier: typedClaims.UserID,
			},
		)
		if err != nil {
			log.Error().Err(err).Msg("could not retrieve user")
			return nil, bizzleException.ErrUnauthorized{Reason: "could not retrieve user: " + err.Error()}
		}

		// create criterion to retrieve user's roles
		roleCriteria := text.List{
			Field: "id",
			List:  make([]string, 0),
		}
		for _, roleID := range findOneUserResponse.User.RoleIDs {
			roleCriteria.List = append(roleCriteria.List, roleID.String())
		}
		roleFindManyResponse, err := a.roleStore.FindMany(
			roleStore.FindManyRequest{
				Criteria: []criterion.Criterion{
					roleCriteria,
					text.Exact{
						Field: "permissions",
						Text:  request.Service,
					},
				},
			},
		)
		if err != nil {
			return nil, bizzleException.ErrUnauthorized{Reason: "could not retrieve roles: " + err.Error()}
		}

		// if any roles match this criteria then the user has access to this service
		if roleFindManyResponse.Total > 0 {
			return &bizzleAuthenticator.AuthenticateServiceResponse{}, nil
		}
		return nil, bizzleException.ErrUnauthorized{Reason: "no permission"}
	}

	return nil, bizzleException.ErrUnauthorized{Reason: "invalid claims"}
}
