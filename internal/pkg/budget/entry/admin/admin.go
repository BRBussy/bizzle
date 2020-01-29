package admin

type Admin interface {
	CreateMany(*CreateManyRequest) (*CreateManyResponse, error)
}

type CreateManyRequest struct {
}

type CreateManyResponse struct {
}
