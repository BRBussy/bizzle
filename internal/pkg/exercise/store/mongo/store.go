package mongo

import (
	"fmt"
	exerciseStore "github.com/BRBussy/bizzle/internal/pkg/exercise/store"
)

type Store struct {
}

func New() exerciseStore.Store {
	return &Store{}
}

func (s Store) Find(request *exerciseStore.FindRequest) (*exerciseStore.FindResponse, error) {
	filter := request.Criteria.ToFilter()
	fmt.Println(filter)
	return &exerciseStore.FindResponse{}, nil
}
