package store

type Store interface {
	Find(request *FindRequest) (*FindResponse, error)
}

type FindRequest struct {
}

type FindResponse struct {
}
