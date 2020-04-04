package basic

import (
	budgetConfig "github.com/BRBussy/bizzle/internal/pkg/budget/config"
	budgetConfigAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/config/admin"
	budgetConfigStore "github.com/BRBussy/bizzle/internal/pkg/budget/config/store"
	budgetConfigValidator "github.com/BRBussy/bizzle/internal/pkg/budget/config/validator"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	validator             validationValidator.Validator
	budgetConfigStore     budgetConfigStore.Store
	budgetConfigValidator budgetConfigValidator.Validator
}

func New(
	validator validationValidator.Validator,
	budgetConfigStore budgetConfigStore.Store,
	budgetConfigValidator budgetConfigValidator.Validator,
) budgetConfigAdmin.Admin {
	return &admin{
		validator:             validator,
		budgetConfigStore:     budgetConfigStore,
		budgetConfigValidator: budgetConfigValidator,
	}
}

func (a *admin) GetMyConfig(request budgetConfigAdmin.GetMyConfigRequest) (*budgetConfigAdmin.GetMyConfigResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	// try and retrieve the budget config for this user
	findOneBudgetConfigResponse, err := a.budgetConfigStore.FindOne(
		budgetConfigStore.FindOneRequest{
			Claims:     request.Claims,
			Identifier: identifier.OwnerID(request.Claims.ScopingID()),
		},
	)
	switch err.(type) {
	case nil:
		// budget config exists for user, return it
		return &budgetConfigAdmin.GetMyConfigResponse{Config: findOneBudgetConfigResponse.Config}, nil

	case mongo.ErrNotFound:
		// budget config doesn't exist yet, create a new blank one and return it
		newConfig := budgetConfig.Config{
			ID:      identifier.ID(uuid.NewV4().String()),
			OwnerID: request.Claims.ScopingID(),
		}
		if _, err := a.budgetConfigStore.CreateOne(
			budgetConfigStore.CreateOneRequest{
				Config: newConfig,
			},
		); err != nil {
			log.Error().Err(err).Msg("could not create budget config")
			return nil, bizzleException.ErrUnexpected{}
		}
		return &budgetConfigAdmin.GetMyConfigResponse{Config: newConfig}, nil

	default:
		// some other retrieval error occurred
		log.Error().Err(err).Msg("could not find budget config")
		return nil, bizzleException.ErrUnexpected{}
	}
}

func (a *admin) SetMyConfig(request budgetConfigAdmin.SetMyConfigRequest) (*budgetConfigAdmin.SetMyConfigResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
}
