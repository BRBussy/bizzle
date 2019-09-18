package identifier

import "strings"

type ErrInvalidIdentifier struct {
	Reasons []string
}

func (e ErrInvalidIdentifier) Error() string {
	return "invalid identifier: " + strings.Join(e.Reasons, ", ")
}
