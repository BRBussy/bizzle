package store

import (
	searchCriterion "github.com/BRBussy/bizzle/pkg/search/criterion"
)

type Store interface {
	Find(request *FindRequest) (*FindResponse, error)
}

const ServiceProvider = "Exercise-Store"
const FindService = ServiceProvider + ".Find"

type FindRequest struct {
	Criteria []searchCriterion.Criterion
}

type FindResponse struct {
}
