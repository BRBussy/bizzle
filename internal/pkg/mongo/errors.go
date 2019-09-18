package mongo

type ErrNotFound struct {
}

func (e ErrNotFound) Error() string {
	return "document not found"
}
