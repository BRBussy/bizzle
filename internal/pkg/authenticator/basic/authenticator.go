package basic

import (
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
)

type Authenticator struct {
}

func (a *Authenticator) Setup() bizzleAuthenticator.Authenticator {
	return &Authenticator{}
}

func (a *Authenticator) Login(*bizzleAuthenticator.LoginRequest) (*bizzleAuthenticator.LoginResponse, error) {
	return &bizzleAuthenticator.LoginResponse{
		JWT: "this has been a success!!!",
	}, nil
}
