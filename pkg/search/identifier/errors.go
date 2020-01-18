package identifier

import "strings"

type ErrInvalidIdentifier struct {
	Reasons []string
}

func (e ErrInvalidIdentifier) Error() string {
	return "invalid identifier: " + strings.Join(e.Reasons, ", ")
}

type ErrInvalidSerializedIdentifier struct {
	Reasons []string
}

func (e ErrInvalidSerializedIdentifier) Error() string {
	return "invalid serialized identifier: " + strings.Join(e.Reasons, ", ")
}

type ErrUnmarshal struct {
	Reasons []string
}

func (e ErrUnmarshal) Error() string {
	return "unmarshalling error: " + strings.Join(e.Reasons, ", ")
}

type ErrMarshal struct {
	Reasons []string
}
