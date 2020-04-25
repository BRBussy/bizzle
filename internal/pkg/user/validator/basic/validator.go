package basic

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	userValidator "github.com/BRBussy/bizzle/internal/pkg/user/validator"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type validator struct {
	requestValidator validationValidator.Validator
	userStore        userStore.Store
	roleStore        roleStore.Store
}

func New(
	requestValidator validationValidator.Validator,
	userStore userStore.Store,
	roleStore roleStore.Store,
) userValidator.Validator {
	return &validator{
		requestValidator: requestValidator,
		userStore:        userStore,
		roleStore:        roleStore,
	}
}

func (v *validator) ValidateForCreate(request userValidator.ValidateForCreateRequest) (*userValidator.ValidateForCreateResponse, error) {
	if err := v.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	reasonsInvalid := make(reasonInvalid.ReasonsInvalid, 0)

	// user must be owned by user invoking service
	if request.User.OwnerID != request.Claims.ScopingID() {
		reasonsInvalid = append(
			reasonsInvalid,
			reasonInvalid.ReasonInvalid{
				Field: "ownerID",
				Type:  reasonInvalid.Invalid,
				Help:  "must be owned by service invoker",
				Data:  request.User.OwnerID,
			},
		)
	}

	// users cannot be registered on creation
	if request.User.Registered {
		reasonsInvalid = append(
			reasonsInvalid,
			reasonInvalid.ReasonInvalid{
				Field: "registered",
				Type:  reasonInvalid.Invalid,
				Help:  "cannot be set during creation",
				Data:  request.User.Registered,
			},
		)
	}

	// each role that a user has assigned needs to exist
	for _, rID := range request.User.RoleIDs {
		if _, err := v.roleStore.FindOne(
			roleStore.FindOneRequest{
				Identifier: rID,
			},
		); err != nil {
			switch err.(type) {
			case mongo.ErrNotFound:
				reasonsInvalid = append(
					reasonsInvalid,
					reasonInvalid.ReasonInvalid{
						Field: "roleIDs",
						Type:  reasonInvalid.DoesntExist,
						Help:  "all roles must exist",
						Data:  rID,
					},
				)
			default:
				log.Error().Err(err).Msg("unable to retrieve user role")
				return nil, bizzleException.ErrUnexpected{}
			}
		}
	}

	// check that user doesn't already exist with this email address
	_, err := v.userStore.FindOne(
		userStore.FindOneRequest{
			Identifier: request.User.Email,
		},
	)
	switch err.(type) {
	case nil:
		// user with this email address already exists
		reasonsInvalid = append(
			reasonsInvalid,
			reasonInvalid.ReasonInvalid{
				Field: "email",
				Type:  reasonInvalid.AlreadyExists,
				Help:  "user with this email address already exists",
				Data:  request.User.Email,
			},
		)

	case mongo.ErrNotFound:
		// do nothing - this is desired
	default:
		log.Error().Err(err).Msg("unable to retrieve user")
		return nil, bizzleException.ErrUnexpected{}
	}

	// password must be blank for creation
	if len(request.User.Password) != 0 {
		reasonsInvalid = append(
			reasonsInvalid,
			reasonInvalid.ReasonInvalid{
				Field: "password",
				Type:  reasonInvalid.Invalid,
				Help:  "password cannot be set during creation",
				Data:  request.User.Email,
			},
		)
	}

	return &userValidator.ValidateForCreateResponse{
		ReasonsInvalid: reasonsInvalid,
	}, nil
}
