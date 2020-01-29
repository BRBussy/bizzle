package basic

import (
	budgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
)

type admin struct {
	validator        validationValidator.Validator
	budgetEntryStore budgetEntryStore.Store
}

func New(
	validator validationValidator.Validator,
	budgetEntryStore budgetEntryStore.Store,
) budgetEntryAdmin.Admin {
	return &admin{
		budgetEntryStore: budgetEntryStore,
		validator:        validator,
	}
}

func (a admin) CreateMany(request *budgetEntryAdmin.CreateManyRequest) (*budgetEntryAdmin.CreateManyResponse, error) {
	if err := a.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &budgetEntryAdmin.CreateManyResponse{}, nil
}
