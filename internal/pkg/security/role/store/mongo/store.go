package mongo

import (
	roleStore "github.com/BRBussy/bizzle/internal/pkg/security/role/store"
)

type store struct {
}

func New() roleStore.Store {
	return &store{}
}

func (s *store) Create(request *roleStore.CreateRequest) (*roleStore.CreateResponse, error) {
	panic("implement me!!!")
}
