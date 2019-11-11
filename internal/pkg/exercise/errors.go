package exercise

import "strings"

type ErrInvalidSerializedExercise struct {
	Reasons []string
}

func (e ErrInvalidSerializedExercise) Error() string {
	return "invalid serialized exercise: " + strings.Join(e.Reasons, ", ")
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
