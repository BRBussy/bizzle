package mongo

import (
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
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
) (userStore.Store, error) {
	// get user collection
	userCollection := database.Collection("user")

	// setup collection indices
	if err := userCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("emailAddress"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up user collection indices")
		return nil, err
	}

	return &store{
		collection: database.Collection("user"),
		validator:  validator,
	}, nil
}

func (s *store) CreateOne(request userStore.CreateOneRequest) (*userStore.CreateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	if err := s.collection.CreateOne(request.User); err != nil {
		log.Error().Err(err).Msg("creating user")
		return nil, err
	}
	return &userStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request userStore.FindOneRequest) (*userStore.FindOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	var result user.User
	if err := s.collection.FindOne(&result, request.Identifier); err != nil {
		switch err.(type) {
		case mongo.ErrNotFound:
			return nil, err
		default:
			log.Error().Err(err).Msg("finding one user")
			return nil, err
		}
	}
	return &userStore.FindOneResponse{User: result}, nil
}

func (s *store) UpdateOne(request userStore.UpdateOneRequest) (*userStore.UpdateOneResponse, error) {
	if err := s.validator.Validate(request); err != nil {
		log.Error().Err(err)
		return nil, err
	}
	if err := s.collection.UpdateOne(request.User, request.User.ID); err != nil {
		log.Error().Err(err).Msg("updating user")
		return nil, err
	}
	return &userStore.UpdateOneResponse{}, nil
}
