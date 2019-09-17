package mongo

import (
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
	"github.com/rs/zerolog/log"
)

type store struct {
	collection *mongo.Collection
}

func New(
	database *mongo.Database,
) roleStore.Store {
	return &store{
		collection: database.Collection("role"),
	}
}

func (s *store) Create(request *roleStore.CreateRequest) (*roleStore.CreateResponse, error) {
	if err := s.collection.CreateOne(request.Role); err != nil {
		log.Error().Err(err).Msg("creating role")
	}
	return &roleStore.CreateResponse{}, nil
}
