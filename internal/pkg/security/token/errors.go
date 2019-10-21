package token

import "strings"

type ErrInvalidToken struct {
	Reasons []string
}

func (e ErrInvalidToken) Error() string {
	return "invalid token: " + strings.Join(e.Reasons, ", ")
}

type ErrTokenVerification struct {
	Reasons []string
}

func (e ErrTokenVerification) Error() string {
	return "token verification failed: " + strings.Join(e.Reasons, ", ")
}
