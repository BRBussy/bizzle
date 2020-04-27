package admin

import "github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"

type ErrUserInvalid struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}

func (e ErrUserInvalid) Error() string {
	return "user invalid: " + e.ReasonsInvalid.String()
}

type ErrUserAlreadyRegistered struct {
}

func (e ErrUserAlreadyRegistered) Error() string {
	return "user already registered"
}

type ErrUserNotRegistered struct {
}

func (e ErrUserNotRegistered) Error() string {
	return "user note registered"
}
