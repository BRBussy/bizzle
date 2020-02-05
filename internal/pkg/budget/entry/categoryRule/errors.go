package categoryRule

type ErrCouldNotClassify struct {
	Reason string
}

func (e ErrCouldNotClassify) Error() string {
	return "could not classify: " + e.Reason
}
