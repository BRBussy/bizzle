package criterion

import "strings"

type ErrInvalid struct {
	Reasons []string
}

func (e ErrInvalid) Error() string {
	return "criterion is invalid: " + strings.Join(e.Reasons, ", ")
}
