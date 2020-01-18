package validator

import "strings"

type ErrRequestNotValid struct {
	Reasons []string
}

func (e ErrRequestNotValid) Error() string {
	return "request not valid: " + strings.Join(e.Reasons, ", ")
}
