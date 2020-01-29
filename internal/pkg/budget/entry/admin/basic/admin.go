package basic

import (
	budgetEntryAdmin "github.com/BRBussy/bizzle/internal/pkg/budget/entry/admin"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
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

	for entryIdx := range request.BudgetEntries {
		request.BudgetEntries[entryIdx].OwnerID = request.Claims.ScopingID()
		request.BudgetEntries[entryIdx].ID = identifier.ID(uuid.NewV4().String())
	}

	if _, err := a.budgetEntryStore.CreateMany(&budgetEntryStore.CreateManyRequest{
		Entries: request.BudgetEntries,
	}); err != nil {
		log.Error().Err(err).Msg("could not create many budget entries")
		return nil, err
	}

	return &budgetEntryAdmin.CreateManyResponse{}, nil
}

func (a admin) DuplicateCheck(*budgetEntryAdmin.DuplicateCheckRequest) (*budgetEntryAdmin.DuplicateCheckResponse, error) {
	panic("implement me")
}
