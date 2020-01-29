package claims

import "strings"

type ErrInvalidSerializedClaims struct {
	Reasons []string
}

func (e ErrInvalidSerializedClaims) Error() string {
	return "invalid serialized claims: " + strings.Join(e.Reasons, ", ")
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

type ErrClaimsNotInContext struct {
}

func (e ErrClaimsNotInContext) Error() string {
	return "claims not in context"
}
