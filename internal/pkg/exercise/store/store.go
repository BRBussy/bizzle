package store

import (
	searchCriteria "github.com/BRBussy/bizzle/pkg/search/criteria"
)

type Store interface {
	Find(request *FindRequest) (*FindResponse, error)
}

const ServiceProvider = "Exercise-Store"
const FindService = ServiceProvider + ".Find"

type FindRequest struct {
	Criteria searchCriteria.Criteria
}

type FindResponse struct {
}
