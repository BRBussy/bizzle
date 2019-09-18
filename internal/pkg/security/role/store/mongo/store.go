package mongo

import (
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/rs/zerolog/log"
	mongoDriver "go.mongodb.org/mongo-driver/mongo"
)

type store struct {
	collection *mongo.Collection
}

func New(
	database *mongo.Database,
) (roleStore.Store, error) {
	// get role collection
	roleCollection := database.Collection("role")

	// setup collection indices
	if err := roleCollection.SetupIndices([]mongoDriver.IndexModel{
		mongo.NewUniqueIndex("id"),
		mongo.NewUniqueIndex("name"),
	}); err != nil {
		log.Error().Err(err).Msg("error setting up role collection indices")
		return nil, err
	}

	return &store{
		collection: database.Collection("role"),
	}, nil
}

func (s *store) Create(request *roleStore.CreateRequest) (*roleStore.CreateResponse, error) {
	if err := s.collection.CreateOne(request.Role); err != nil {
		log.Error().Err(err).Msg("creating role")
	}
	return &roleStore.CreateResponse{}, nil
}

func (s *store) FindOne(request *roleStore.FindOneRequest) (*roleStore.FindOneResponse, error) {

}
