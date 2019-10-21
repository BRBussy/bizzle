package basic

import (
	bizzleAuthenticator "github.com/BRBussy/bizzle/internal/pkg/authenticator"
	userStore "github.com/BRBussy/bizzle/internal/pkg/user/store"
	"github.com/BRBussy/bizzle/pkg/search/identifier"
	"golang.org/x/crypto/bcrypt"
)

type Authenticator struct {
	userStore userStore.Store
}

func (a *Authenticator) Setup(
	userStore userStore.Store,
) bizzleAuthenticator.Authenticator {
	return &Authenticator{
		userStore: userStore,
	}
}

func (a *Authenticator) Login(request *bizzleAuthenticator.LoginRequest) (*bizzleAuthenticator.LoginResponse, error) {
	// try and retrieve user by email address
	retrieveResponse, err := a.userStore.FindOne(&userStore.FindOneRequest{
		Identifier: identifier.Email(request.Email),
	})
	if err != nil {
		return nil, err
	}

	// check password is correct
	if err := bcrypt.CompareHashAndPassword(retrieveResponse.User.Password, []byte(request.Password)); err != nil {
		return nil, err
	}

	// generate login claims

	return &bizzleAuthenticator.LoginResponse{
		JWT: "this has been a success!!!",
	}, nil
}
