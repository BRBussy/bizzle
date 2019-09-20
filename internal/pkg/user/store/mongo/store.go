package mongo

import (
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	"github.com/BRBussy/bizzle/internal/pkg/user"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	collection *mongo.Collection
}

func New(
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
	}, nil
}

func (s *store) CreateOne(request *userStore.CreateOneRequest) (*userStore.CreateOneResponse, error) {
	if err := s.collection.CreateOne(request.User); err != nil {
		log.Error().Err(err).Msg("creating user")
		return nil, err
	}
	return &userStore.CreateOneResponse{}, nil
}

func (s *store) FindOne(request *userStore.FindOneRequest) (*userStore.FindOneResponse, error) {
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

func (s *store) UpdateOne(request *userStore.UpdateOneRequest) (*userStore.UpdateOneResponse, error) {
	if err := s.collection.UpdateOne(request.User, identifier.ID(request.User.ID)); err != nil {
		log.Error().Err(err).Msg("updating user")
		return nil, err
	}
	return &userStore.UpdateOneResponse{}, nil
}
