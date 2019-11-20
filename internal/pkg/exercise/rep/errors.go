package rep

import "strings"

type ErrInvalidSerializedRep struct {
	Reasons []string
}

func (e ErrInvalidSerializedRep) Error() string {
	return "invalid serialized rep: " + strings.Join(e.Reasons, ", ")
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
