package criterion

import "strings"

type ErrInvalid struct {
	Reasons []string
}

func (e ErrInvalid) Error() string {
	return "criterion is invalid: " + strings.Join(e.Reasons, ", ")
}

type ErrUnexpected struct {
	Reasons []string
}

func (e ErrUnexpected) Error() string {
	return "unexpected error: " + strings.Join(e.Reasons, ", ")
}
