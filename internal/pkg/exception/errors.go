package exception

import "strings"

type ErrUnexpected struct {
	Reasons []string
}

func (e ErrUnexpected) Error() string {
	return "unexpected error: " + strings.Join(e.Reasons, ", ")
}

type ErrUnauthorized struct {
	Reason string
}

func (e ErrUnauthorized) Error() string {
	if e.Reason == "" {
		return "unauthorized"
	}
	return "unauthorized " + e.Reason
}
