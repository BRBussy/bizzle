package mongo

import (
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
)

type Store struct {
}

func New() exerciseStore.Store {
	return &Store{}
}

func (s Store) Find(request *exerciseStore.FindRequest) (*exerciseStore.FindResponse, error) {
	return &exerciseStore.FindResponse{}, nil
}
