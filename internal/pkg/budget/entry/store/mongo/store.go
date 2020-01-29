package mongo

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
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
) (budgetEntryStore.Store, error) {
	// get budgetEntry collection
	budgetEntryCollection := database.Collection("budgetEntry")

	// setup collection indices
	if err := budgetEntryCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("name", "variant"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up budgetEntry collection indices")
		return nil, err
	}

	return &store{
		validator:  validator,
		collection: budgetEntryCollection,
	}, nil
}

func (s store) CreateOne(request *budgetEntryStore.CreateOneRequest) (*budgetEntryStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	if err := s.collection.CreateOne(request.Entry); err != nil {
		log.Error().Err(err).Msg("creating role")
		return nil, err
	}
	return &budgetEntryStore.CreateOneResponse{}, nil
}

func (s store) FindOne(request *budgetEntryStore.FindOneRequest) (*budgetEntryStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var result budgetEntry.Entry
	if err := s.collection.FindOne(&result, request.Identifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one budgetEntry")
			return nil, err
		}
	}

	return &budgetEntryStore.FindOneResponse{
		Entry: result,
	}, nil
}

func (s store) FindMany(request *budgetEntryStore.FindManyRequest) (*budgetEntryStore.FindManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var records []budgetEntry.Entry
	count, err := s.collection.FindMany(&records, request.Criteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding exercises")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]budgetEntry.Entry, 0)
	}

	return &budgetEntryStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s store) UpdateOne(request *budgetEntryStore.UpdateOneRequest) (*budgetEntryStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.UpdateOne(request.Entry, request.Entry.ID); err != nil {
		log.Error().Err(err).Msg("updating budgetEntry")
		return nil, err
	}

	return &budgetEntryStore.UpdateOneResponse{}, nil
}
