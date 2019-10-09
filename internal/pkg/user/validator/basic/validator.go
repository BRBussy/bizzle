package basic

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userValidator "github.com/BRBussy/bizzle/internal/pkg/user/validator"
	"github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"
	"github.com/rs/zerolog/log"
)

type validator struct {
	roleStore roleStore.Store
}

func (v validator) genericValidation(userToValidate user.User) ([]reasonInvalid.ReasonInvalid, error) {
	reasonsInvalid := make([]reasonInvalid.ReasonInvalid, 0)

	if userToValidate.Name == "" {
		reasonsInvalid = append(reasonsInvalid,
			reasonInvalid.ReasonInvalid{
				Field: "name",
				Type:  reasonInvalid.Blank,
				Help:  "can't be blank",
				Data:  userToValidate.Name,
			})
	}

	if userToValidate.Email == "" {
		reasonsInvalid = append(reasonsInvalid,
			reasonInvalid.ReasonInvalid{
				Field: "email",
				Type:  reasonInvalid.Blank,
				Help:  "can't be blank",
				Data:  userToValidate.Email,
			})
	}

	// roles cannot be nil
	if userToValidate.RoleIDs == nil {
		reasonsInvalid = append(reasonsInvalid,
			reasonInvalid.ReasonInvalid{
				Field: "roleIDs",
				Type:  reasonInvalid.Nil,
				Help:  "can't be nil",
				Data:  userToValidate.RoleIDs,
			})
	} else {
		// all roles must exist
		for i := range userToValidate.RoleIDs {
			if _, err := v.roleStore.FindOne(&roleStore.FindOneRequest{
				Identifier: userToValidate.RoleIDs[i],
			}); err != nil {
				switch err.(type) {
				case mongo.ErrNotFound:
					reasonsInvalid = append(reasonsInvalid,
						reasonInvalid.ReasonInvalid{
							Field: "roleIDs",
							Type:  reasonInvalid.DoesntExist,
							Help:  "must exist",
							Data:  userToValidate.RoleIDs[i],
						})
				default:
					log.Error().Err(err).Msg("retrieving role")
					return nil, bizzleException.ErrUnexpected{}
				}
			}
		}
	}

	return reasonsInvalid, nil
}

func (v validator) ValidateForCreate(request *userValidator.ValidateForCreateRequest) (*userValidator.ValidateForCreateResponse, error) {
	reasonsInvalid, err := v.genericValidation(request.User)
	if err != nil {
		log.Error().Err(err).Msg("generic validation")
		return nil, err
	}

	return &userValidator.ValidateForCreateResponse{ReasonsInvalid: reasonsInvalid}, nil
}

func New(
	roleStore roleStore.Store,
) userValidator.Validator {
	return &validator{
		roleStore: roleStore,
	}
}
