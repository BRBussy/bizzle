package store

type Store interface {
	Find(request *FindRequest) (*FindResponse, error)
}

const ServiceProvider = "Exercise-Store"
const FindService = ServiceProvider + ".Find"

type FindRequest struct {
}

type FindResponse struct {
}
