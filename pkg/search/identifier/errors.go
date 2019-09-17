package identifier

import "strings"

type ErrInvalid struct {
	Reasons []string
}

func (e ErrInvalid) Error() string {
	return "invalid identifier: " + strings.Join(e.Reasons, ", ")
}
