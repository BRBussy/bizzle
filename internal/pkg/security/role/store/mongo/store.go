package mongo

import (
	"github.com/BRBussy/bizzle/internal/pkg/mongo"
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
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
	panic("implement me!!!")
}
