package basic

import (
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
)

type authenticator struct {
}

func (a *authenticator) Setup() bizzleAuthenticator.Authenticator {
	return &authenticator{}
}

func (a *authenticator) Login(*bizzleAuthenticator.LoginRequest) (*bizzleAuthenticator.LoginResponse, error) {
	return &bizzleAuthenticator.LoginResponse{
		JWT: "this has been a success!!!",
	}, nil
}
