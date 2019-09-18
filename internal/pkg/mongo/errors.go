package mongo

import "strings"

type ErrNotFound struct {
}

func (e ErrNotFound) Error() string {
	return "document not found"
}

type ErrInvalidIdentifier struct {
	Reasons []string
}

func (e ErrInvalidIdentifier) Error() string {
	return "invalid identifier: " + strings.Join(e.Reasons, ", ")
}
