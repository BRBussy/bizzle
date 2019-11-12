package basic

import (
	exerciseAdmin "github.com/BRBussy/bizzle/internal/pkg/exercise/admin"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type admin struct {
	validator     validationValidator.Validator
	exerciseStore exerciseStore.Store
}

func New(
	validator validationValidator.Validator,
	exerciseStore exerciseStore.Store,
) exerciseAdmin.Admin {
	return &admin{
		exerciseStore: exerciseStore,
		validator:     validator,
	}
}

func (a admin) CreateOne(request *exerciseAdmin.CreateOneRequest) (*exerciseAdmin.CreateOneResponse, error) {
	if err := a.validator.ValidateRequest(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	request.Exercise.SetID(identifier.ID(uuid.NewV4().String()))

	if _, err := a.exerciseStore.CreateOne(&exerciseStore.CreateOneRequest{
		Exercise: request.Exercise,
	}); err != nil {
		log.Error().Err(err).Msg("creating exercise")
		return nil, err
	}

	return &exerciseAdmin.CreateOneResponse{Exercise: request.Exercise}, nil
}

func (a admin) UpdateOne(request *exerciseAdmin.UpdateOneRequest) (*exerciseAdmin.UpdateOneResponse, error) {
	if err := a.validator.ValidateRequest(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	return &exerciseAdmin.UpdateOneResponse{}, nil
}
