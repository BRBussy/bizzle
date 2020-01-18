package admin

import "github.com/BRBussy/bizzle/pkg/validate/reasonInvalid"

type ErrUserInvalid struct {
	ReasonsInvalid reasonInvalid.ReasonsInvalid
}

func (e ErrUserInvalid) Error() string {
	return "user invalid: " + e.ReasonsInvalid.String()
}
