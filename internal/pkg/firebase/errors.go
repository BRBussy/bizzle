package firebase

type ErrUnexpected struct {
}

func (e ErrUnexpected) Error() string {
	return "unexpected firebase error"
}
