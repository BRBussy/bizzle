package mongo

import (
	budgetEntry "github.com/BRBussy/bizzle/internal/pkg/budget/entry"
	budgetEntryStore "github.com/BRBussy/bizzle/internal/pkg/budget/entry/store"
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	validationValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
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
	}); err != nil {
		log.Error().Err(err).Msg("error setting up budgetEntry collection indices")
		return nil, err
	}

	return &store{
		validator:  validator,
		collection: budgetEntryCollection,
	}, nil
}

func (s *store) CreateOne(request *budgetEntryStore.CreateOneRequest) (*budgetEntryStore.CreateOneResponse, error) {
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

func (s *store) CreateMany(request *budgetEntryStore.CreateManyRequest) (*budgetEntryStore.CreateManyResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	documentsToCreate := make([]interface{}, 0)
	for _, entry := range request.Entries {
		documentsToCreate = append(documentsToCreate, entry)
	}
	if err := s.collection.CreateMany(documentsToCreate); err != nil {
		log.Error().Err(err).Msg("creating budget entries")
		return nil, err
	}

	return &budgetEntryStore.CreateManyResponse{}, nil
}

func (s *store) FindOne(request *budgetEntryStore.FindOneRequest) (*budgetEntryStore.FindOneResponse, error) {
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

func (s *store) FindMany(request *budgetEntryStore.FindManyRequest) (*budgetEntryStore.FindManyResponse, error) {
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

func (s *store) FindManyComposite(request *budgetEntryStore.FindManyCompositeRequest) (*budgetEntryStore.FindManyCompositeResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var result []budgetEntry.CompositeEntry
	pipeline := []bson.D{
		{
			{
				Key: "$lookup",
				Value: bson.M{
					"from":         "budgetEntryCategoryRule",
					"localField":   "categoryRuleID",
					"foreignField": "id",
					"as":           "categoryRule",
				},
			},
		},

		// $unwind stage
		// for every categoryRule entity in the entry.categoryRule array created in previous lookup
		// stage a new document will be output with the category rule entity at user.person
		// (this assumes that there is only 1 category rule per entry)
		// result:
		//
		// {
		// 	id: 1234,
		//  ...otherUserFields,
		//  categoryRule: { id: 4321, ...otherPersonFields},
		// }
		{
			{
				Key:   "$unwind",
				Value: "$categoryRule",
			},
		},

		// $match Stage
		// apply given filter
		{
			{
				Key:   "$match",
				Value: request.Criteria.ToFilter(),
			},
		},
	}

	count, err := s.collection.Aggregate(
		pipeline,
		request.Query,
		&result,
	)
	if err != nil {
		log.Error().Err(err).Msg("could not find users")
		return nil, err
	}
	if result == nil {
		result = make([]budgetEntry.CompositeEntry, 0)
	}

	return &budgetEntryStore.FindManyCompositeResponse{
		Records: result,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request *budgetEntryStore.UpdateOneRequest) (*budgetEntryStore.UpdateOneResponse, error) {
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
