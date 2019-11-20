package mongo

import (
	bizzleException "github.com/BRBussy/bizzle/internal/pkg/exception"
	"github.com/BRBussy/bizzle/internal/pkg/exercise/session"
	sessionStore "github.com/BRBussy/bizzle/internal/pkg/exercise/session/store"
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	validateValidator "github.com/BRBussy/bizzle/pkg/validate/validator"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	collection       *mongo.Collection
	requestValidator validateValidator.Validator
}

func New(
	requestValidator validateValidator.Validator,
	database *mongo.Database,
) (sessionStore.Store, error) {
	// get session collection
	sessionCollection := database.Collection("session")

	// setup collection indices
	if err := sessionCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up session collection indices")
		return nil, err
	}

	return &store{
		collection:       database.Collection("session"),
		requestValidator: requestValidator,
	}, nil
}

func (s *store) CreateOne(request *sessionStore.CreateOneRequest) (*sessionStore.CreateOneResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.CreateOne(request.Session); err != nil {
		log.Error().Err(err).Msg("creating session")
		return nil, err
	}
	return &sessionStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request *sessionStore.FindOneRequest) (*sessionStore.FindOneResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var result session.Session
	if err := s.collection.FindOne(&result, request.Identifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one session")
			return nil, err
		}
	}
	return &sessionStore.FindOneResponse{Session: result}, nil
}

func (s *store) FindMany(request *sessionStore.FindManyRequest) (*sessionStore.FindManyResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	var records []session.Session
	count, err := s.collection.FindMany(&records, request.Criteria, request.Query)
	if err != nil {
		log.Error().Err(err).Msg("finding sessions")
		return nil, bizzleException.ErrUnexpected{}
	}
	if records == nil {
		records = make([]session.Session, 0)
	}

	return &sessionStore.FindManyResponse{
		Records: records,
		Total:   count,
	}, nil
}

func (s *store) UpdateOne(request *sessionStore.UpdateOneRequest) (*sessionStore.UpdateOneResponse, error) {
	if err := s.requestValidator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}

	if err := s.collection.UpdateOne(request.Session, identifier.ID(request.Session.ID)); err != nil {
		log.Error().Err(err).Msg("updating session")
		return nil, err
	}
	return &sessionStore.UpdateOneResponse{}, nil
}
