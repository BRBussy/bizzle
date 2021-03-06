package mongo

import "strings"

type ErrInvalidConfig struct {
	Reasons []string
}

func (e ErrInvalidConfig) Error() string {
	return "invalid config: " + strings.Join(e.Reasons, ", ")
}

type ErrUnexpected struct {
}

func (e ErrUnexpected) Error() string {
	return "unexpected mongo error"
}

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

type ErrInvalidCriteria struct {
	Reasons []string
}

func (e ErrInvalidCriteria) Error() string {
	return "invalid criteria: " + strings.Join(e.Reasons, ", ")
}

type ErrQueryInvalid struct {
	Reasons []string
}

func (e ErrQueryInvalid) Error() string {
	return "invalid query: " + strings.Join(e.Reasons, ", ")
}

type ErrSortingInvalid struct {
	Reasons []string
}

func (e ErrSortingInvalid) Error() string {
	return "sorting invalid: " + strings.Join(e.Reasons, ", ")
}
