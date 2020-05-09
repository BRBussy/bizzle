package basic

import (
	budgetEntryIgnoredAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored/admin"
	budgetEntryIgnoredStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/ignored/store"
	"github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	validator               validationValidator.Validator
	budgetEntryIgnoredStore budgetEntryIgnoredStore.Store
}

func New(
	validator validationValidator.Validator,
	budgetEntryIgnoredStore budgetEntryIgnoredStore.Store,
) budgetEntryIgnoredAdmin.Admin {
	return &admin{
		validator:               validator,
		budgetEntryIgnoredStore: budgetEntryIgnoredStore,
	}
}

func (a admin) CreateOne(request budgetEntryIgnoredAdmin.CreateOneRequest) (*budgetEntryIgnoredAdmin.CreateOneResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	request.Ignored.ID = identifier.ID(uuid.NewV4().String())
	request.Ignored.OwnerID = request.Claims.ScopingID()

	if _, err := a.budgetEntryIgnoredStore.CreateOne(
		budgetEntryIgnoredStore.CreateOneRequest{
			Ignored: request.Ignored,
		},
	); err != nil {
		log.Error().Err(err).Msg("unable to create ignored")
		return nil, exception.ErrUnexpected{}
	}

	return &budgetEntryIgnoredAdmin.CreateOneResponse{}, nil
}
