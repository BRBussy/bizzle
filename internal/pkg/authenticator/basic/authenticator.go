package basic

import (
	brizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
)

type authenticator struct {
}

func New() brizzleAuthenticator.Authenticator {
	return &authenticator{}
}

func (a *authenticator) SignUp(*brizzleAuthenticator.SignUpRequest) (*brizzleAuthenticator.SignUpResponse, error) {
	return &brizzleAuthenticator.SignUpResponse{
		Msg: "this has been a success!!!",
	}, nil
}
