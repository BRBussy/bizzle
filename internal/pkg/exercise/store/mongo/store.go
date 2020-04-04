package mongo

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/exercise"
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

func (s *store) CreateOne(request exerciseStore.CreateOneRequest) (*exerciseStore.CreateOneResponse, error) {
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

func (s *store) FindOne(request exerciseStore.FindOneRequest) (*exerciseStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var result exercise.Exercise
	if err := s.collection.FindOne(&result, request.Identifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one exercise")
			return nil, err
		}
	}

	return &exerciseStore.FindOneResponse{
		Exercise: result,
	}, nil
}

func (s *store) FindMany(request exerciseStore.FindManyRequest) (*exerciseStore.FindManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var records []exercise.Exercise
	count, err := s.collection.FindMany(&records, request.Criteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding exercises")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]exercise.Exercise, 0)
	}

	return &exerciseStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request exerciseStore.UpdateOneRequest) (*exerciseStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.UpdateOne(request.Exercise, request.Exercise.ID); err != nil {
		log.Error().Err(err).Msg("updating exercise")
		return nil, err
	}

	return &exerciseStore.UpdateOneResponse{}, nil
}
