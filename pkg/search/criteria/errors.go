package criteria

import "strings"

type ErrInvalidSerializedCriteria struct {
	Reasons []string
}

func (e ErrInvalidSerializedCriteria) Error() string {
	return "serialized criteria is invalid: " + strings.Join(e.Reasons, ", ")
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

func (e ErrMarshal) Error() string {
	return "marshalling error: " + strings.Join(e.Reasons, ", ")
}
