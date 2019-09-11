package exception

import "strings"

type ErrUnexpected struct {
	Reasons []string
}

func (e ErrUnexpected) Error() string {
	return "unexpected error: " + strings.Join(e.Reasons, ", ")
}
