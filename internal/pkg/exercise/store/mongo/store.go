package mongo

import (
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	validator  validationValidator.Validator
	collection *mongo.Collection
}

func New(
	validator validationValidator.Validator,
	database *mongo.Database,
) (exerciseStore.Store, error) {
	// get exercise collection
	exerciseCollection := database.Collection("exercise")

	// setup collection indices
	if err := exerciseCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("name", "variant"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up exercise collection indices")
		return nil, err
	}

	return &store{
		validator:  validator,
		collection: exerciseCollection,
	}, nil
}

func (s store) CreateOne(request *exerciseStore.CreateOneRequest) (*exerciseStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	if err := s.collection.CreateOne(request.Exercise); err != nil {
		log.Error().Err(err).Msg("creating role")
		return nil, err
	}
	return &exerciseStore.CreateOneResponse{}, nil
}

func (s store) FindOne(request *exerciseStore.FindOneRequest) (*exerciseStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return &exerciseStore.FindOneResponse{}, nil
}

func (s store) FindMany(request *exerciseStore.FindManyRequest) (*exerciseStore.FindManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return &exerciseStore.FindManyResponse{}, nil
}

func (s store) UpdateOne(request *exerciseStore.UpdateOneRequest) (*exerciseStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return &exerciseStore.UpdateOneResponse{}, nil
}
